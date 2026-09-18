import ballerina/io;
import ballerina/test;

function squareDataSet() returns map<[int, int]> {
    return {"one": [1, 1], "two": [2, 4], "three": [3, 9]};
}

function testSquare(int input, int expected) {
    test:assertEquals(input * input, expected);
}

public function main() {
    test:setTestOptions("target", "ddmod", "ddmod", "false", "false", "", "", "", "false", "false");
    test:registerTestConfig("testSquare", testSquare, true, [], [], (), (), false, squareDataSet);
    int exitCode = test:startSuite();
    io:println("exitCode: ", exitCode);
    // @output 		[pass] testSquare#one
    // @output 		[pass] testSquare#two
    // @output 		[pass] testSquare#three
    // @output 
    // @output 
    // @output 		3 passing
    // @output 		0 failing
    // @output 		0 skipped
    // @output 
    // @output 		Test execution time : <DURATION>s
    // @output exitCode: 0
}
