# Context
You need to migrate a toml parser from java to golang. The java toml parser implementation is in /Users/asmaj/ballerina/ballerina-lang/misc/toml-parser. The golang equivalent should be implemented in /Users/asmaj/ballerina/ballerina-lang-go/toml-parser.

## instructions
Understand the architecture  of the java implementation, package structure, APIs and anything needed to implement in golang. 
Document the architecture.
List down all the java classes
Create a 1:1 mapping for java classes to golang implementation
Create plan for the migration and track progress

## Best practices
- Use names of java equivalent  as much as possible
- Create separate files in go to match classes in java
- Add comments inline to show the java equivalent 
- Use the java implementation as source of truth
- Ask for approval when for conflicting instructions
- Use golang best practices
- After every phase, validate the implementation against the architecture