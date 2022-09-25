***Important**: The IAC and Go code is for research purposes only and not recommended to be used in a production environment*

# Enterprise Observability
Enterprise observability refers to an organization's ability to measure the state of their enterprise based on the external output of their systems.

## Gaining Visibility into a 3-tier Architecture
This repo provides Terraform IAC files to create a simple 3-tier architecture on Azure. The accompanying Go code then queries Azure Resource Graph APIs to upload entities and their relationships to a Neo4J graph database.

### Step 1: Cloning the Repository
Run the following command to install the repository locally.

`go install https://github.com/jose-canul/enterprise-observability@latest`

### Step 2: Create the test environment
Once in the project directory of your terminal, move into the deployments/iac directory.

`cd deployments/iac`

Run terraform to build the test environment and select yes to create the test infrastructure.

`terraform apply`

### Step 3: Run the observability app.
In the project directory, run the command below. Please make sure that you have a Neo4J instance up and connected to the application.

`go run main.go`