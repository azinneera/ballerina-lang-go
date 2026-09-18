import ballerina/io;
import ballerina/test;

int oneRan = 0;
int twoRan = 0;

function squareDataSet() returns map<[int, int]> {
    return {"one": [1, 1], "two": [2, 4]};
}

function testSquare(int input, int expected) {
    if input == 1 {
        oneRan += 1;
    } else {
        twoRan += 1;
    }
    test:assertEquals(input * input, expected);
}

public function main() {
    // --tests testSquare#one selects only the "one" sub-case of a map-shaped
    // data provider — exercises skipDataDrivenTest's sub-key matching, a
    // distinct code path from plain by-name --tests filtering.
    test:setTestOptions("target", "ddsubkeymod", "ddsubkeymod", "false", "false", "", "", "testSquare#one",
            "false", "false");
    test:registerTestConfig("testSquare", testSquare, true, [], [], (), (), false, squareDataSet);

    int exitCode = test:startSuite();

    io:println("exitCode: ", exitCode);
    io:println("oneRan: ", oneRan);
    io:println("twoRan: ", twoRan);
    // @output 		[pass] testSquare#one
    // @output 
    // @output 
    // @output 		1 passing
    // @output 		0 failing
    // @output 		0 skipped
    // @output 
    // @output 		Test execution time : <DURATION>s
    // @output exitCode: 0
    // @output oneRan: 1
    // @output twoRan: 0
}
