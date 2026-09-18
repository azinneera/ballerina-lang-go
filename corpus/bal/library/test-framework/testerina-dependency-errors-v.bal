import ballerina/io;
import ballerina/test;

function testA() {
}

function testB() {
}

public function main() {
    // Cyclic dependency: testA depends on testB, testB depends on testA.
    test:setTestOptions("target", "cyclemod", "cyclemod", "false", "false", "", "", "", "false", "false");
    test:registerTestConfig("testA", testA, true, [], [testB], (), (), false, ());
    test:registerTestConfig("testB", testB, true, [], [testA], (), (), false, ());
    int exitCode = test:startSuite();
    io:println("cyclic exitCode: ", exitCode);
    // @output Cyclic test dependencies detected: testA -> testB -> testA
    // @output cyclic exitCode: 1
}
