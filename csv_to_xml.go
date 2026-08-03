package main

import (
	"encoding/csv"
	"encoding/xml"
	"flag"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

var xmlNameRegexp = regexp.MustCompile(`[^A-Za-z0-9_:-]`)

func main() {
	inputPath := flag.String("in", "Accounts_50000.csv", "input CSV file path")
	outputPath := flag.String("out", "Accounts_50000.xml", "output XML file path")
	rootName := flag.String("root", "Accounts", "root XML element name")
	rowName := flag.String("row", "Account", "row XML element name")
	flag.Parse()

	if *outputPath == "" {
		base := strings.TrimSuffix(filepath.Base(*inputPath), filepath.Ext(*inputPath))
		*outputPath = fmt.Sprintf("%s.xml", base)
	}

	if err := convertCSVToXML(*inputPath, *outputPath, *rootName, *rowName); err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}
	fmt.Printf("Converted %s -> %s\n", *inputPath, *outputPath)
}

func convertCSVToXML(csvPath, xmlPath, rootName, rowName string) error {
	csvFile, err := os.Open(csvPath)
	if err != nil {
		return fmt.Errorf("open CSV file: %w", err)
	}
	defer csvFile.Close()

	reader := csv.NewReader(csvFile)
	headers, err := reader.Read()
	if err != nil {
		return fmt.Errorf("read CSV header: %w", err)
	}
	if len(headers) == 0 {
		return fmt.Errorf("CSV header is empty")
	}

	for i := range headers {
		headers[i] = sanitizeXMLName(headers[i])
		if headers[i] == "" {
			headers[i] = fmt.Sprintf("Field%d", i+1)
		}
	}

	xmlFile, err := os.Create(xmlPath)
	if err != nil {
		return fmt.Errorf("create XML file: %w", err)
	}
	defer xmlFile.Close()

	xmlFile.WriteString(xml.Header)
	encoder := xml.NewEncoder(xmlFile)
	encoder.Indent("", "  ")

	root := xml.StartElement{Name: xml.Name{Local: sanitizeXMLName(rootName)}}
	row := xml.StartElement{Name: xml.Name{Local: sanitizeXMLName(rowName)}}

	if err := encoder.EncodeToken(root); err != nil {
		return fmt.Errorf("write XML root start: %w", err)
	}

	for {
		record, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return fmt.Errorf("read CSV record: %w", err)
		}
		if len(record) != len(headers) {
			return fmt.Errorf("unexpected record length: got %d values, want %d", len(record), len(headers))
		}

		if err := encoder.EncodeToken(row); err != nil {
			return fmt.Errorf("write row start: %w", err)
		}
		for i, value := range record {
			field := xml.StartElement{Name: xml.Name{Local: headers[i]}}
			if err := encoder.EncodeElement(value, field); err != nil {
				return fmt.Errorf("write field %q: %w", headers[i], err)
			}
		}
		if err := encoder.EncodeToken(row.End()); err != nil {
			return fmt.Errorf("write row end: %w", err)
		}
	}

	if err := encoder.EncodeToken(root.End()); err != nil {
		return fmt.Errorf("write XML root end: %w", err)
	}
	if err := encoder.Flush(); err != nil {
		return fmt.Errorf("flush XML writer: %w", err)
	}

	return nil
}

func sanitizeXMLName(name string) string {
	name = strings.TrimSpace(name)
	name = xmlNameRegexp.ReplaceAllString(name, "_")
	if name == "" {
		return ""
	}
	if !isNameStartChar(name[0]) {
		name = "_" + name
	}
	return name
}

func isNameStartChar(b byte) bool {
	return (b >= 'A' && b <= 'Z') || (b >= 'a' && b <= 'z') || b == '_' || b == ':'
}
