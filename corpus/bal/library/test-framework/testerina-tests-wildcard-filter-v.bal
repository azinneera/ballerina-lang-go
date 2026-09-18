import ballerina/io;
import ballerina/test;

int testFooRan = 0;
int testBarRan = 0;
int otherBazRan = 0;

function testFoo() {
    testFooRan += 1;
}

function testBar() {
    testBarRan += 1;
}

function otherBaz() {
    otherBazRan += 1;
}

public function main() {
    // --tests "test*" selects every function whose name matches the
    // wildcard pattern (filter.bal#hasTest -> matchWildcard), not just exact
    // names — testFoo/testBar match, otherBaz doesn't.
    test:setTestOptions("target", "wildcardfiltermod", "wildcardfiltermod", "false", "false", "", "", "test*",
            "false", "false");
    test:registerTestConfig("testFoo", testFoo, true, [], [], (), (), false, ());
    test:registerTestConfig("testBar", testBar, true, [], [], (), (), false, ());
    test:registerTestConfig("otherBaz", otherBaz, true, [], [], (), (), false, ());

    int exitCode = test:startSuite();

    io:println("exitCode: ", exitCode);
    io:println("testFooRan: ", testFooRan);
    io:println("testBarRan: ", testBarRan);
    io:println("otherBazRan: ", otherBazRan);
    // @output 		[pass] testBar
    // @output 		[pass] testFoo
    // @output
    // @output
    // @output 		2 passing
    // @output 		0 failing
    // @output 		0 skipped
    // @output
    // @output 		Test execution time : <DURATION>s
    // @output exitCode: 0
    // @output testFooRan: 1
    // @output testBarRan: 1
    // @output otherBazRan: 0
}
