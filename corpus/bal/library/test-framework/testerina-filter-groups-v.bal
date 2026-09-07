import ballerina/io;
import ballerina/test;

int unitRan = 0;
int integrationRan = 0;

function testUnit() {
    unitRan += 1;
}

function testIntegration() {
    integrationRan += 1;
}

public function main() {
    // setTestOptions must run before registration: registerTestConfig's
    // group-filter gating reads filterGroups at registration time, not
    // lazily (see TODO.md/CHECKPOINT.md's P8.6-P8.8 entry).
    test:setTestOptions("target", "filtermod", "filtermod", "false", "false", "unit", "", "", "false", "false");
    test:registerTestConfig("testUnit", testUnit, true, ["unit"], [], (), (), false, ());
    test:registerTestConfig("testIntegration", testIntegration, true, ["integration"], [], (), (), false, ());

    int exitCode = test:startSuite();

    io:println("exitCode: ", exitCode);
    io:println("unitRan: ", unitRan);
    io:println("integrationRan: ", integrationRan);
    // @output 		[pass] testUnit
    // @output 
    // @output 
    // @output 		1 passing
    // @output 		0 failing
    // @output 		0 skipped
    // @output 
    // @output 		Test execution time : <DURATION>s
    // @output exitCode: 0
    // @output unitRan: 1
    // @output integrationRan: 0
}
