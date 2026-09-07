import ballerina/io;
import ballerina/test;

function testDisableFunc1() {
    test:assertTrue(false, "this test is not expected to run");
}

function testDisableFunc2() {
    test:assertTrue(true);
}

function testDisableFunc3() {
    test:assertFalse(false);
}

public function main() {
    // A disabled (enable: false) test never runs and isn't counted at all —
    // not even as skipped. testDisableFunc2 (enable: true) and
    // testDisableFunc3 (no enable field, defaults to true) both run
    // normally.
    test:setTestOptions("target", "disabledtestsmod", "disabledtestsmod", "false", "false", "", "", "",
            "false", "false");
    test:registerTestConfig("testDisableFunc1", testDisableFunc1, false, [], [], (), (), false, ());
    test:registerTestConfig("testDisableFunc2", testDisableFunc2, true, [], [], (), (), false, ());
    test:registerTestConfig("testDisableFunc3", testDisableFunc3, true, [], [], (), (), false, ());

    int exitCode = test:startSuite();

    io:println("exitCode: ", exitCode);
    // @output 		[pass] testDisableFunc2
    // @output 		[pass] testDisableFunc3
    // @output
    // @output
    // @output 		2 passing
    // @output 		0 failing
    // @output 		0 skipped
    // @output
    // @output 		Test execution time : <DURATION>s
    // @output exitCode: 0
}
