// SPDX-License-Identifier: Apache-2.0

package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/mendixlabs/mxcli/mdl/visitor"
	"github.com/spf13/cobra"
)

var describeCmd = &cobra.Command{
	Use:   "describe <type> <name>",
	Short: "Describe a project element",
	Long: `Describe an element from a Mendix project in MDL syntax.

Types:
  entity           Describe an entity
  externalentity   Describe an external entity (alias for entity)
  association      Describe an association
  enumeration      Describe an enumeration
  constant         Describe a constant
  microflow        Describe a microflow
  workflow         Describe a workflow
  page             Describe a page
  snippet          Describe a snippet
  layout           Describe a layout
  javaaction       Describe a java action
  odataclient      Describe a consumed OData service
  odataservice     Describe a published OData service
  businesseventservice  Describe a business event service (also: "business event service")
  modulerole       Describe a module role
  userrole         Describe a user role
  projectsecurity  Show project security settings
  demouser         Describe a demo user
  navigation       Describe a navigation profile
  systemoverview   Module dependency graph (requires --format elk)

Example:
  mxcli describe -p app.mpr entity MyModule.Customer
  mxcli describe -p app.mpr microflow MyModule.ProcessOrder
  mxcli describe -p app.mpr workflow MyModule.MyWorkflow
  mxcli describe -p app.mpr navigation Responsive
  mxcli describe -p app.mpr layout Atlas_Core.PopupLayout
  mxcli describe -p app.mpr javaaction MyModule.MyAction
  mxcli describe -p app.mpr odataclient MyModule.ExternalAPI
  mxcli describe -p app.mpr odataservice MyModule.CustomerAPI
  mxcli describe -p app.mpr businesseventservice MyModule.CustomerEventsApi
  mxcli describe -p app.mpr business event service MyModule.CustomerEventsApi
  mxcli describe -p app.mpr constant MyModule.BaseUrl
  mxcli describe -p app.mpr modulerole MyModule.User
  mxcli describe -p app.mpr userrole Administrator
  mxcli describe -p app.mpr projectsecurity ProjectSecurity
  mxcli describe -p app.mpr --format elk systemoverview SystemOverview
`,
	Args: cobra.MinimumNArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		projectPath, _ := cmd.Flags().GetString("project")
		if projectPath == "" {
			fmt.Fprintln(os.Stderr, "Error: --project (-p) is required")
			os.Exit(1)
		}

		// Support multi-word types: "business event service Module.Name" → type="BUSINESS EVENT SERVICE", name="Module.Name"
		// The last arg is always the name, everything before is the type.
		name := args[len(args)-1]
		objectType := strings.ToUpper(strings.Join(args[:len(args)-1], " "))

		var mdlCmd string
		switch objectType {
		case "ENTITY":
			mdlCmd = fmt.Sprintf("DESCRIBE ENTITY %s", name)
		case "ASSOCIATION":
			mdlCmd = fmt.Sprintf("DESCRIBE ASSOCIATION %s", name)
		case "ENUMERATION":
			mdlCmd = fmt.Sprintf("DESCRIBE ENUMERATION %s", name)
		case "MICROFLOW":
			mdlCmd = fmt.Sprintf("DESCRIBE MICROFLOW %s", name)
		case "WORKFLOW":
			mdlCmd = fmt.Sprintf("DESCRIBE WORKFLOW %s", name)
		case "PAGE":
			mdlCmd = fmt.Sprintf("DESCRIBE PAGE %s", name)
		case "SNIPPET":
			mdlCmd = fmt.Sprintf("DESCRIBE SNIPPET %s", name)
		case "LAYOUT":
			mdlCmd = fmt.Sprintf("DESCRIBE LAYOUT %s", name)
		case "MODULEROLE", "MODULE ROLE":
			mdlCmd = fmt.Sprintf("DESCRIBE MODULE ROLE %s", name)
		case "USERROLE", "USER ROLE":
			mdlCmd = fmt.Sprintf("DESCRIBE USER ROLE '%s'", name)
		case "PROJECTSECURITY", "PROJECT SECURITY":
			mdlCmd = "SHOW PROJECT SECURITY"
		case "DEMOUSER", "DEMO USER":
			mdlCmd = fmt.Sprintf("DESCRIBE DEMO USER '%s'", name)
		case "JAVAACTION", "JAVA ACTION":
			mdlCmd = fmt.Sprintf("DESCRIBE JAVA ACTION %s", name)
		case "CONSTANT":
			mdlCmd = fmt.Sprintf("DESCRIBE CONSTANT %s", name)
		case "ODATACLIENT", "ODATA CLIENT":
			mdlCmd = fmt.Sprintf("DESCRIBE ODATA CLIENT %s", name)
		case "ODATASERVICE", "ODATA SERVICE":
			mdlCmd = fmt.Sprintf("DESCRIBE ODATA SERVICE %s", name)
		case "BUSINESSEVENTSERVICE", "BUSINESS EVENT SERVICE":
			mdlCmd = fmt.Sprintf("DESCRIBE BUSINESS EVENT SERVICE %s", name)
		case "DATABASECONNECTION", "DATABASE CONNECTION":
			mdlCmd = fmt.Sprintf("DESCRIBE DATABASE CONNECTION %s", name)
		case "EXTERNALENTITY", "EXTERNAL ENTITY":
			mdlCmd = fmt.Sprintf("DESCRIBE ENTITY %s", name)
		case "NAVIGATION":
			mdlCmd = fmt.Sprintf("DESCRIBE NAVIGATION %s", name)
		case "NAVPROFILE":
			mdlCmd = fmt.Sprintf("DESCRIBE NAVIGATION %s", name)
		case "SYSTEMOVERVIEW":
			mdlCmd = "" // handled directly by format-specific path
		default:
			fmt.Fprintf(os.Stderr, "Unknown type: %s\n", strings.Join(args[:len(args)-1], " "))
			fmt.Fprintln(os.Stderr, "Valid types: entity, association, enumeration, constant, microflow, workflow, page, snippet, layout, javaaction, odataclient, odataservice, businesseventservice, databaseconnection, modulerole, userrole, projectsecurity, demouser, navigation, systemoverview")
			fmt.Fprintln(os.Stderr, "Multi-word types also accepted: business event service, odata service, java action, database connection, etc.")
			os.Exit(1)
		}

		exec, logger := newLoggedExecutor("subcommand")
		defer logger.Close()
		defer exec.Close()
		exec.SetQuiet(true) // suppress status messages for programmatic output

		// Connect
		connectProg, _ := visitor.Build(fmt.Sprintf("CONNECT LOCAL '%s'", projectPath))
		for _, stmt := range connectProg.Statements {
			if err := exec.Execute(stmt); err != nil {
				fmt.Fprintf(os.Stderr, "Error: %v\n", err)
				os.Exit(1)
			}
		}

		// Check for mermaid/elk format - bypass MDL parser and call directly
		format, _ := cmd.Flags().GetString("format")
		typeArg := strings.Join(args[:len(args)-1], " ")
		if format == "mermaid" {
			if err := exec.DescribeMermaid(typeArg, name); err != nil {
				fmt.Fprintf(os.Stderr, "Error: %v\n", err)
				os.Exit(1)
			}
			return
		} else if format == "elk" {
			upper := objectType
			if upper == "SYSTEMOVERVIEW" {
				if err := exec.ModuleOverview(); err != nil {
					fmt.Fprintf(os.Stderr, "Error: %v\n", err)
					os.Exit(1)
				}
			} else if upper == "ENTITY" || upper == "DOMAINMODEL" || upper == "EXTERNALENTITY" || upper == "EXTERNAL ENTITY" {
				if err := exec.DomainModelELK(name); err != nil {
					fmt.Fprintf(os.Stderr, "Error: %v\n", err)
					os.Exit(1)
				}
			} else if upper == "MICROFLOW" {
				if err := exec.MicroflowELK(name); err != nil {
					fmt.Fprintf(os.Stderr, "Error: %v\n", err)
					os.Exit(1)
				}
			} else if upper == "PAGE" {
				if err := exec.PageWireframeJSON(name); err != nil {
					fmt.Fprintf(os.Stderr, "Error: %v\n", err)
					os.Exit(1)
				}
			} else if upper == "SNIPPET" {
				if err := exec.SnippetWireframeJSON(name); err != nil {
					fmt.Fprintf(os.Stderr, "Error: %v\n", err)
					os.Exit(1)
				}
			} else {
				fmt.Fprintf(os.Stderr, "ELK format not supported for type: %s\n", typeArg)
				os.Exit(1)
			}
			return
		}

		// SYSTEMOVERVIEW requires elk format
		if mdlCmd == "" {
			fmt.Fprintf(os.Stderr, "Type %s requires --format elk\n", args[0])
			os.Exit(1)
		}

		// Execute describe command
		descProg, errs := visitor.Build(mdlCmd)
		if len(errs) > 0 {
			for _, err := range errs {
				fmt.Fprintf(os.Stderr, "Parse error: %v\n", err)
			}
			os.Exit(1)
		}

		for _, stmt := range descProg.Statements {
			if err := exec.Execute(stmt); err != nil {
				fmt.Fprintf(os.Stderr, "Error: %v\n", err)
				os.Exit(1)
			}
		}
	},
}
