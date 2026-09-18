import ballerina/io;
import ballerina/test;

int unitRan = 0;
int fastRan = 0;
int slowRan = 0;

function testUnit() {
    unitRan += 1;
}

function testFast() {
    fastRan += 1;
}

function testSlow() {
    slowRan += 1;
}

public function main() {
    // --groups "unit,fast" (comma-separated) selects a test if it belongs
    // to ANY of the listed groups, not just one — distinct from
    // testerina-filter-groups-v.bal, which only exercises a single group.
    test:setTestOptions("target", "multigroupmod", "multigroupmod", "false", "false", "unit,fast", "", "",
            "false", "false");
    test:registerTestConfig("testUnit", testUnit, true, ["unit"], [], (), (), false, ());
    test:registerTestConfig("testFast", testFast, true, ["fast"], [], (), (), false, ());
    test:registerTestConfig("testSlow", testSlow, true, ["slow"], [], (), (), false, ());

    int exitCode = test:startSuite();

    io:println("exitCode: ", exitCode);
    io:println("unitRan: ", unitRan);
    io:println("fastRan: ", fastRan);
    io:println("slowRan: ", slowRan);
    // @output 		[pass] testFast
    // @output 		[pass] testUnit
    // @output
    // @output
    // @output 		2 passing
    // @output 		0 failing
    // @output 		0 skipped
    // @output
    // @output 		Test execution time : <DURATION>s
    // @output exitCode: 0
    // @output unitRan: 1
    // @output fastRan: 1
    // @output slowRan: 0
}
