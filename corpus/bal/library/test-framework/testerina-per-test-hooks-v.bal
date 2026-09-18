import ballerina/io;
import ballerina/test;

string[] callOrder = [];

function myBefore() {
    callOrder.push("before");
}

function myAfter() {
    callOrder.push("after");
}

function myTest() {
    callOrder.push("test");
    test:assertTrue(true);
}

public function main() {
    test:setTestOptions("target", "pertesthooksmod", "pertesthooksmod", "false", "false", "", "", "", "false", "false");
    // Per-test before:/after: (distinct from global beforeEach/afterEach) —
    // both must run exactly once, in before -> test -> after order.
    test:registerTestConfig("myTest", myTest, true, [], [], myBefore, myAfter, false, ());

    int exitCode = test:startSuite();

    io:println("exitCode: ", exitCode);
    io:println("callOrder: ", callOrder);
    // @output 		[pass] myTest
    // @output 
    // @output 
    // @output 		1 passing
    // @output 		0 failing
    // @output 		0 skipped
    // @output 
    // @output 		Test execution time : <DURATION>s
    // @output exitCode: 0
    // @output callOrder: ["before","test","after"]
}
