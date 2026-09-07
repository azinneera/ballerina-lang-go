import ballerina/io;
import ballerina/test;

int matchRan = 0;
int otherRan = 0;

function codeDataSet() returns map<[string]> {
    return {
        "a,b\"c": ["match"],
        "other": ["other"]
    };
}

function testFunction3(string value1) {
    if value1 == "match" {
        matchRan += 1;
    } else {
        otherRan += 1;
    }
}

public function main() {
    // --tests testFunction3#a%2Cb%22c selects the sub-case keyed by the
    // literal string `a,b"c` — a comma and a double-quote, both special
    // characters that must be percent-encoded on the CLI (filter.bal's
    // skipDataDrivenTest round-trips this via escapeSpecialCharacters, the
    // same mechanism jballerina's own StringUtils#escapeSpecialCharacters
    // uses for its "code fragment" sub-keys).
    test:setTestOptions("target", "specialcharkeymod", "specialcharkeymod", "false", "false", "", "",
            "testFunction3#a%2Cb%22c", "false", "false");
    test:registerTestConfig("testFunction3", testFunction3, true, [], [], (), (), false, codeDataSet);

    int exitCode = test:startSuite();

    io:println("exitCode: ", exitCode);
    io:println("matchRan: ", matchRan);
    io:println("otherRan: ", otherRan);
    // @output 		[pass] testFunction3#a,b"c
    // @output
    // @output
    // @output 		1 passing
    // @output 		0 failing
    // @output 		0 skipped
    // @output
    // @output 		Test execution time : <DURATION>s
    // @output exitCode: 0
    // @output matchRan: 1
    // @output otherRan: 0
}
