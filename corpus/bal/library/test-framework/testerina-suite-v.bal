import ballerina/io;
import ballerina/test;

int beforeSuiteCount = 0;
int afterSuiteCount = 0;
int beforeEachCount = 0;
int afterEachCount = 0;
string[] executionOrder = [];

function beforeAll() {
    beforeSuiteCount += 1;
}

function afterAll() {
    afterSuiteCount += 1;
}

function beforeEachTest() {
    beforeEachCount += 1;
}

function afterEachTest() {
    afterEachCount += 1;
}

function testSetup() {
    executionOrder.push("testSetup");
}

function testMainCase() {
    executionOrder.push("testMainCase");
    test:assertTrue(true);
}

// A failing test is deliberately not included here: ballerina/test's own
// failure-reporting (formatFailedError) unconditionally appends a trailing
// tabs-only line to any failure message (lines.push(""); join with tabs),
// which this corpus harness's @output annotations can never represent
// (outputRe's parser trims trailing space/tab from every annotated line —
// confirmed via a standalone repro). Failure-path output is already covered
// by cli/cmd/test_test.go's Go-level string comparisons, which have no such
// trimming. This test focuses on the parts @output *can* verify exactly:
// hook ordering and counts, and dependsOn execution order, all-passing.
public function main() {
    test:setTestOptions("target", "suitemod", "suitemod", "false", "false", "", "", "", "false", "false");
    test:registerBeforeSuite("beforeAll", beforeAll);
    test:registerAfterSuite("afterAll", afterAll, true);
    test:registerBeforeEach("beforeEachTest", beforeEachTest);
    test:registerAfterEach("afterEachTest", afterEachTest);
    test:registerTestConfig("testSetup", testSetup, true, [], [], (), (), false, ());
    test:registerTestConfig("testMainCase", testMainCase, true, [], [testSetup], (), (), false, ());

    int exitCode = test:startSuite();

    io:println("exitCode: ", exitCode);
    io:println("beforeSuiteCount: ", beforeSuiteCount);
    io:println("afterSuiteCount: ", afterSuiteCount);
    io:println("beforeEachCount: ", beforeEachCount);
    io:println("afterEachCount: ", afterEachCount);
    io:println("executionOrder: ", executionOrder);
    // @output 		[pass] testSetup
    // @output 		[pass] testMainCase
    // @output 
    // @output 
    // @output 		2 passing
    // @output 		0 failing
    // @output 		0 skipped
    // @output 
    // @output 		Test execution time : <DURATION>s
    // @output exitCode: 0
    // @output beforeSuiteCount: 1
    // @output afterSuiteCount: 1
    // @output beforeEachCount: 2
    // @output afterEachCount: 2
    // @output executionOrder: ["testSetup","testMainCase"]
}
