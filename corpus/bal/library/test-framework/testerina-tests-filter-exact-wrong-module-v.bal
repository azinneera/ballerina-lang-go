import ballerina/io;
import ballerina/test;

int fooTestRan = 0;
int barTestRan = 0;

function fooTest() {
    fooTestRan += 1;
}

function barTest() {
    barTestRan += 1;
}

public function main() {
    // "combinedqualmod.sub:fooTest" is an exact filter entry naming this
    // module's fooTest, but qualified to a *different* module within the
    // same package — it must not shadow "combinedqualmod:foo*", a wildcard
    // entry that also matches fooTest and does target the running module.
    // filter.bal#hasTest must fall through to checking the wildcard filters
    // even when an exact entry exists for the same name but the wrong
    // module — otherwise fooTest is wrongly excluded entirely, never
    // reaching the wildcard that would select it.
    test:setTestOptions("target", "combinedqualmod", "combinedqualmod", "false", "false", "", "",
            "combinedqualmod.sub:fooTest,combinedqualmod:foo*", "false", "false");
    test:registerTestConfig("fooTest", fooTest, true, [], [], (), (), false, ());
    test:registerTestConfig("barTest", barTest, true, [], [], (), (), false, ());

    int exitCode = test:startSuite();

    io:println("exitCode: ", exitCode);
    io:println("fooTestRan: ", fooTestRan);
    io:println("barTestRan: ", barTestRan);
    // @output 		[pass] fooTest
    // @output
    // @output
    // @output 		1 passing
    // @output 		0 failing
    // @output 		0 skipped
    // @output
    // @output 		Test execution time : <DURATION>s
    // @output exitCode: 0
    // @output fooTestRan: 1
    // @output barTestRan: 0
}
