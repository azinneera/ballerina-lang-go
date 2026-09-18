import ballerina/io;
import ballerina/test;

int testFuncRan = 0;
int testFunc2Ran = 0;

function testFunc() {
    testFuncRan += 1;
}

function testFunc2() {
    testFunc2Ran += 1;
}

public function main() {
    // --tests testFunc2 selects only the dependent function, but its
    // dependency (testFunc) is automatically pulled in and run too — a
    // dependency-incomplete --tests selection isn't an error, jballerina
    // and this port both auto-include the dependency chain.
    test:setTestOptions("target", "testsfilterdepmod", "testsfilterdepmod", "false", "false", "", "",
            "testFunc2", "false", "false");
    test:registerTestConfig("testFunc", testFunc, true, [], [], (), (), false, ());
    test:registerTestConfig("testFunc2", testFunc2, true, [], [testFunc], (), (), false, ());

    int exitCode = test:startSuite();

    io:println("exitCode: ", exitCode);
    io:println("testFuncRan: ", testFuncRan);
    io:println("testFunc2Ran: ", testFunc2Ran);
    // @output 		[pass] testFunc
    // @output 		[pass] testFunc2
    // @output
    // @output
    // @output 		2 passing
    // @output 		0 failing
    // @output 		0 skipped
    // @output
    // @output 		Test execution time : <DURATION>s
    // @output exitCode: 0
    // @output testFuncRan: 1
    // @output testFunc2Ran: 1
}
