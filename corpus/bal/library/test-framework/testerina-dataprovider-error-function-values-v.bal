import ballerina/io;
import ballerina/test;

function errorDataSet() returns map<[error, string]> {
    error e = error("foo");
    return {"foo": [e, "foo"]};
}

function getMessage(error input) returns string {
    return input.message();
}

function testErrorData(error input, string expected) {
    test:assertEquals(getMessage(input), expected);
}

function isEven(int x) returns boolean {
    return x % 2 == 0;
}

function isPositive(int x) returns boolean {
    return x > 0;
}

function functionDataSet() returns map<[function, int]> {
    return {
        "even": [isEven, 4],
        "positive": [isPositive, 2]
    };
}

function testFunctionData(function (int x) returns boolean fn, int value) {
    test:assertTrue(fn(value));
}

public function main() {
    // Data providers aren't limited to primitive/record-shaped values —
    // an `error` value and a `function` value can both flow through as
    // ordinary test data, dynamically invoked/inspected inside the test
    // function like any other argument.
    test:setTestOptions("target", "dataprovidererrorfnmod", "dataprovidererrorfnmod", "false", "false", "", "",
            "", "false", "false");
    test:registerTestConfig("testErrorData", testErrorData, true, [], [], (), (), false, errorDataSet);
    test:registerTestConfig("testFunctionData", testFunctionData, true, [], [], (), (), false, functionDataSet);

    int exitCode = test:startSuite();

    io:println("exitCode: ", exitCode);
    // @output 		[pass] testErrorData#foo
    // @output 		[pass] testFunctionData#even
    // @output 		[pass] testFunctionData#positive
    // @output
    // @output
    // @output 		3 passing
    // @output 		0 failing
    // @output 		0 skipped
    // @output
    // @output 		Test execution time : <DURATION>s
    // @output exitCode: 0
}
