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
    // --disable-groups "unit,fast" (comma-separated) excludes a test if it
    // belongs to ANY of the listed groups — distinct from
    // testerina-disable-groups-v.bal, which only exercises a single
    // disabled group.
    test:setTestOptions("target", "multidisablegroupsmod", "multidisablegroupsmod", "false", "false", "",
            "unit,fast", "", "false", "false");
    test:registerTestConfig("testUnit", testUnit, true, ["unit"], [], (), (), false, ());
    test:registerTestConfig("testFast", testFast, true, ["fast"], [], (), (), false, ());
    test:registerTestConfig("testSlow", testSlow, true, ["slow"], [], (), (), false, ());

    int exitCode = test:startSuite();

    io:println("exitCode: ", exitCode);
    io:println("unitRan: ", unitRan);
    io:println("fastRan: ", fastRan);
    io:println("slowRan: ", slowRan);
    // @output 		[pass] testSlow
    // @output
    // @output
    // @output 		1 passing
    // @output 		0 failing
    // @output 		0 skipped
    // @output
    // @output 		Test execution time : <DURATION>s
    // @output exitCode: 0
    // @output unitRan: 0
    // @output fastRan: 0
    // @output slowRan: 1
}
