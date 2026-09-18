import ballerina/io;
import ballerina/test;

int unitRan = 0;
int slowRan = 0;
int fastRan = 0;

function testUnit() {
    unitRan += 1;
}

function testSlow() {
    slowRan += 1;
}

function testFast() {
    fastRan += 1;
}

public function main() {
    // --disable-groups slow excludes only the "slow"-grouped test; "unit"
    // and "fast" (unrelated groups) must still run.
    test:setTestOptions("target", "disablegroupsmod", "disablegroupsmod", "false", "false", "", "slow", "", "false", "false");
    test:registerTestConfig("testUnit", testUnit, true, ["unit"], [], (), (), false, ());
    test:registerTestConfig("testSlow", testSlow, true, ["slow"], [], (), (), false, ());
    test:registerTestConfig("testFast", testFast, true, ["fast"], [], (), (), false, ());

    int exitCode = test:startSuite();

    io:println("exitCode: ", exitCode);
    io:println("unitRan: ", unitRan);
    io:println("slowRan: ", slowRan);
    io:println("fastRan: ", fastRan);
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
    // @output slowRan: 0
    // @output fastRan: 1
}
