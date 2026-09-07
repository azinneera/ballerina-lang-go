import ballerina/io;
import ballerina/test;

int beforeUnitGroupCount = 0;
int afterUnitGroupCount = 0;

function beforeUnitGroup() {
    beforeUnitGroupCount += 1;
}

function afterUnitGroup() {
    afterUnitGroupCount += 1;
}

function testUnitOne() {
    test:assertTrue(true);
}

function testUnitTwo() {
    test:assertTrue(true);
}

public function main() {
    test:setTestOptions("target", "groupshooksmod", "groupshooksmod", "false", "false", "", "", "", "false", "false");
    test:registerBeforeGroups("beforeUnitGroup", beforeUnitGroup, ["unit"]);
    test:registerAfterGroups("afterUnitGroup", afterUnitGroup, ["unit"], false);
    test:registerTestConfig("testUnitOne", testUnitOne, true, ["unit"], [], (), (), false, ());
    test:registerTestConfig("testUnitTwo", testUnitTwo, true, ["unit"], [], (), (), false, ());

    int exitCode = test:startSuite();

    io:println("exitCode: ", exitCode);
    // beforeGroups/afterGroups run once per group, not once per test.
    io:println("beforeUnitGroupCount: ", beforeUnitGroupCount);
    io:println("afterUnitGroupCount: ", afterUnitGroupCount);
    // @output 		[pass] testUnitOne
    // @output 		[pass] testUnitTwo
    // @output 
    // @output 
    // @output 		2 passing
    // @output 		0 failing
    // @output 		0 skipped
    // @output 
    // @output 		Test execution time : <DURATION>s
    // @output exitCode: 0
    // @output beforeUnitGroupCount: 1
    // @output afterUnitGroupCount: 1
}
