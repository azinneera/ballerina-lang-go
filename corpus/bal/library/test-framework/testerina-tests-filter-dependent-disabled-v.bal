import ballerina/io;
import ballerina/test;

function testDisabledFunc() {
    test:assertTrue(false, "should not run");
}

function testDependentDisabledFunc() {
}

public function main() {
    // --tests testDependentDisabledFunc selects a function that depends on
    // a disabled one — same "depends on function ... disabled or not
    // included" error as testerina-disabled-dependency-v.bal, just reached
    // via --tests selection instead of a plain full run.
    test:setTestOptions("target", "testsfilterdepdisabledmod", "testsfilterdepdisabledmod", "false", "false",
            "", "", "testDependentDisabledFunc", "false", "false");
    test:registerTestConfig("testDisabledFunc", testDisabledFunc, false, [], [], (), (), false, ());
    test:registerTestConfig("testDependentDisabledFunc", testDependentDisabledFunc, true, [], [testDisabledFunc],
            (), (), false, ());

    int exitCode = test:startSuite();

    io:println("exitCode: ", exitCode);
    // @output error: Test [testDependentDisabledFunc] depends on function [testDisabledFunc], but it is either disabled or not included.
    // @output exitCode: 1
}
