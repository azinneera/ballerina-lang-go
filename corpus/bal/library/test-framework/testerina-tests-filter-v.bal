import ballerina/io;
import ballerina/test;

int testAlphaRan = 0;
int testBetaRan = 0;
int testGammaRan = 0;

function testAlpha() {
    testAlphaRan += 1;
}

function testBeta() {
    testBetaRan += 1;
}

function testGamma() {
    testGammaRan += 1;
}

public function main() {
    // --tests testAlpha,testBeta selects two of three by name; a
    // non-existing name in the list is simply ignored (no error), matching
    // jballerina's own behavior.
    test:setTestOptions("target", "testsfiltermod", "testsfiltermod", "false", "false", "", "",
            "testAlpha,testBeta,noSuchTest", "false", "false");
    test:registerTestConfig("testAlpha", testAlpha, true, [], [], (), (), false, ());
    test:registerTestConfig("testBeta", testBeta, true, [], [], (), (), false, ());
    test:registerTestConfig("testGamma", testGamma, true, [], [], (), (), false, ());

    int exitCode = test:startSuite();

    io:println("exitCode: ", exitCode);
    io:println("testAlphaRan: ", testAlphaRan);
    io:println("testBetaRan: ", testBetaRan);
    io:println("testGammaRan: ", testGammaRan);
    // @output 		[pass] testAlpha
    // @output 		[pass] testBeta
    // @output 
    // @output 
    // @output 		2 passing
    // @output 		0 failing
    // @output 		0 skipped
    // @output 
    // @output 		Test execution time : <DURATION>s
    // @output exitCode: 0
    // @output testAlphaRan: 1
    // @output testBetaRan: 1
    // @output testGammaRan: 0
}
