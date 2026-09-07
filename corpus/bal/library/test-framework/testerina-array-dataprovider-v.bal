import ballerina/io;
import ballerina/test;

// Tuple/array-shaped data provider (AnyOrError[][]) — distinct from
// testerina-datadriven-v.bal's map-shaped (map<[int,int]>) provider. Sub-test
// suffixes are the tuple's own index (0, 1, 2, ...), not a caller-chosen key.
function squareDataSet() returns [int, int][] {
    return [[1, 1], [2, 4], [3, 9]];
}

function testSquare(int input, int expected) {
    test:assertEquals(input * input, expected);
}

public function main() {
    test:setTestOptions("target", "arrayddmod", "arrayddmod", "false", "false", "", "", "", "false", "false");
    test:registerTestConfig("testSquare", testSquare, true, [], [], (), (), false, squareDataSet);

    int exitCode = test:startSuite();

    io:println("exitCode: ", exitCode);
    // @output 		[pass] testSquare#0
    // @output 		[pass] testSquare#1
    // @output 		[pass] testSquare#2
    // @output 
    // @output 
    // @output 		3 passing
    // @output 		0 failing
    // @output 		0 skipped
    // @output 
    // @output 		Test execution time : <DURATION>s
    // @output exitCode: 0
}
