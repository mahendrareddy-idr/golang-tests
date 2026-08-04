package main

import (
	"encoding/csv"
	"fmt"
	"log"
	"math/rand"
	"os"
	"strconv"
	"time"
)

func main() {
	rand.Seed(time.Now().UnixNano())
	numRecords := 50000

	// --- Generate Accounts ---
	accountFile, err := os.Create("Accounts_50000.csv")
	if err != nil {
		log.Fatal(err)
	}
	defer accountFile.Close()

	accountWriter := csv.NewWriter(accountFile)
	defer accountWriter.Flush()

	accountsHeader := []string{"AccountName", "BillingStreet", "BillingCity", "BillingState", "BillingPostalCode", "BillingCountry", "Phone", "Industry"}
	accountWriter.Write(accountsHeader)

	accountNames := make([]string, numRecords)
	industries := []string{"Technology", "Finance", "Healthcare", "Manufacturing"}

	for i := 0; i < numRecords; i++ {
		accountName := fmt.Sprintf("Account_%d", i+1)
		accountNames[i] = accountName
		record := []string{
			accountName,
			fmt.Sprintf("Street %d", i+1),
			fmt.Sprintf("City %d", i+1),
			"CA",
			fmt.Sprintf("%05d", 10000+rand.Intn(90000)),
			"USA",
			fmt.Sprintf("555-%03d-%04d", rand.Intn(999), rand.Intn(9999)),
			industries[rand.Intn(len(industries))],
		}
		accountWriter.Write(record)
	}

	// --- Generate Leads ---
	leadsFile, err := os.Create("Leads_50000.csv")
	if err != nil {
		log.Fatal(err)
	}
	defer leadsFile.Close()

	leadsWriter := csv.NewWriter(leadsFile)
	defer leadsWriter.Flush()

	leadsHeader := []string{"FirstName", "LastName", "Company", "Email", "Phone", "LeadSource", "Status"}
	leadsWriter.Write(leadsHeader)

	leadSources := []string{"Web", "Referral", "Partner"}
	statuses := []string{"Open - Not Contacted", "Working - Contacted", "Closed - Converted"}

	for i := 0; i < numRecords; i++ {
		record := []string{
			fmt.Sprintf("First%d", i+1),
			fmt.Sprintf("Last%d", i+1),
			fmt.Sprintf("Company_%d", i+1),
			fmt.Sprintf("user%d@example.com", i+1),
			fmt.Sprintf("555-%03d-%04d", rand.Intn(999), rand.Intn(9999)),
			leadSources[rand.Intn(len(leadSources))],
			statuses[rand.Intn(len(statuses))],
		}
		leadsWriter.Write(record)
	}

	// --- Generate Contacts (1-to-1 with Accounts) ---
	contactsFile, err := os.Create("Contacts_50000.csv")
	if err != nil {
		log.Fatal(err)
	}
	defer contactsFile.Close()

	contactsWriter := csv.NewWriter(contactsFile)
	defer contactsWriter.Flush()

	contactsHeader := []string{"FirstName", "LastName", "Email", "Phone", "AccountName"}
	contactsWriter.Write(contactsHeader)

	for i := 0; i < numRecords; i++ {
		record := []string{
			fmt.Sprintf("ContactFirst%d", i+1),
			fmt.Sprintf("ContactLast%d", i+1),
			fmt.Sprintf("contact%d@example.com", i+1),
			fmt.Sprintf("555-%03d-%04d", rand.Intn(999), rand.Intn(9999)),
			accountNames[i], // 1-to-1 mapping to Account
		}
		contactsWriter.Write(record)
	}

	// --- Generate Opportunities (1-to-1 with Accounts) ---
	opportunitiesFile, err := os.Create("Opportunities_50000.csv")
	if err != nil {
		log.Fatal(err)
	}
	defer opportunitiesFile.Close()

	opportunitiesWriter := csv.NewWriter(opportunitiesFile)
	defer opportunitiesWriter.Flush()

	opportunitiesHeader := []string{"OpportunityName", "AccountName", "StageName", "CloseDate", "Amount"}
	opportunitiesWriter.Write(opportunitiesHeader)

	stages := []string{"Qualification", "Proposal/Price Quote", "Negotiation/Review", "Closed Won", "Closed Lost"}

	for i := 0; i < numRecords; i++ {
		closeDate := time.Now().AddDate(0, rand.Intn(12), rand.Intn(28)).Format("2006-01-02")
		amount := strconv.Itoa(rand.Intn(99001) + 1000) // 1000 - 100000
		record := []string{
			fmt.Sprintf("Opportunity_%d", i+1),
			accountNames[i], // 1-to-1 mapping to Account
			stages[rand.Intn(len(stages))],
			closeDate,
			amount,
		}
		opportunitiesWriter.Write(record)
	}

	fmt.Println("CSV files generated: Accounts_50000.csv, Leads_50000.csv, Contacts_50000.csv, Opportunities_50000.csv")
}
