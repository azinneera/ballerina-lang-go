import ballerina/io;
import ballerina/test;

function testDisabledFunc() {
    test:assertTrue(false, "should not run");
}

public function main() {
    // Explicitly selecting a disabled (enable: false) function by name via
    // --tests doesn't override the disable flag — it's simply not
    // included, same as if it had never been named at all.
    test:setTestOptions("target", "testsfilterdisabledmod", "testsfilterdisabledmod", "false", "false", "", "",
            "testDisabledFunc", "false", "false");
    test:registerTestConfig("testDisabledFunc", testDisabledFunc, false, [], [], (), (), false, ());

    int exitCode = test:startSuite();

    io:println("exitCode: ", exitCode);
    // @output
    // @output
    // @output 		No tests found
    // @output
    // @output 		Test execution time : <DURATION>s
    // @output exitCode: 0
}
