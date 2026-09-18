import ballerina/io;
import ballerina/test;

function emptyDataSet() returns map<[int, int]> {
    return {};
}

function testSquare(int input, int expected) {
    test:assertEquals(input * input, expected);
}

public function main() {
    // A data provider that's validly-shaped but returns zero entries — the
    // test function itself is registered, but zero sub-cases ever get
    // expanded/executed, so startSuite() reports "No tests found" rather
    // than "0 passing / 0 failing / 0 skipped". Confirmed non-crashing,
    // sane behavior — distinct from the tracked type/arity-MISMATCH bug
    // (#21), which is about a data provider's return shape not matching
    // its test function's actual parameters, not an empty-but-valid one.
    test:setTestOptions("target", "emptydpmod", "emptydpmod", "false", "false", "", "", "", "false", "false");
    test:registerTestConfig("testSquare", testSquare, true, [], [], (), (), false, emptyDataSet);

    int exitCode = test:startSuite();

    io:println("exitCode: ", exitCode);
    // @output
    // @output
    // @output 		No tests found
    // @output
    // @output 		Test execution time : <DURATION>s
    // @output exitCode: 0
}
