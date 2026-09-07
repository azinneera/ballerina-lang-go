import ballerina/io;
import ballerina/test;

function testA() {
}

function testB() {
}

public function main() {
    // testB depends on testA, but testA is explicitly disabled.
    test:setTestOptions("target", "disabledmod", "disabledmod", "false", "false", "", "", "", "false", "false");
    test:registerTestConfig("testA", testA, false, [], [], (), (), false, ());
    test:registerTestConfig("testB", testB, true, [], [testA], (), (), false, ());
    int exitCode = test:startSuite();
    io:println("disabled-dependency exitCode: ", exitCode);
    // @output error: Test [testB] depends on function [testA], but it is either disabled or not included.
    // @output disabled-dependency exitCode: 1
}
