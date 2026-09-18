import ballerina/test;

function testOne() {
}

public function main() {
    test:setTestOptions("target", "listgroupsmod", "listgroupsmod", "false", "false", "", "", "", "false", "true");
    test:registerTestConfig("testOne", testOne, true, ["unit", "fast"], [], (), (), false, ());
    int exitCode = test:startSuite();
    _ = exitCode;
    // @output 	[unit, fast]
}
