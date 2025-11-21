import assert from "node:assert/strict";
import { afterEach, describe, it } from "node:test";
import { getArch, getPlatform } from "../../../src/utils/platforms";
import { mockArch, mockPlatform } from "../../helpers/test-utils";

describe("Platform Utilities (Improved)", () => {
  describe("getArch()", () => {
    let restore: (() => void) | undefined;

    afterEach(() => {
      restore?.();
    });

    it("should map arm64 to aarch64", () => {
      restore = mockArch("arm64");
      const result = getArch();
      assert.equal(result, "aarch64");
    });

    it("should map ia32 to i686", () => {
      restore = mockArch("ia32");
      const result = getArch();
      assert.equal(result, "i686");
    });

    it("should map x64 to x86_64", () => {
      restore = mockArch("x64");
      const result = getArch();
      assert.equal(result, "x86_64");
    });

    it("should return undefined for unsupported architecture", () => {
      restore = mockArch("unsupported");
      const result = getArch();
      assert.equal(result, undefined);
    });
  });

  describe("getPlatform()", () => {
    let restore: (() => void) | undefined;

    afterEach(() => {
      restore?.();
    });

    it("should map darwin to apple-darwin", () => {
      restore = mockPlatform("darwin");
      const result = getPlatform();
      assert.equal(result, "apple-darwin");
    });

    it("should map linux to unknown-linux-gnu", () => {
      restore = mockPlatform("linux");
      const result = getPlatform();
      assert.equal(result, "unknown-linux-gnu");
    });

    it("should map win32 to pc-windows-msvc", () => {
      restore = mockPlatform("win32");
      const result = getPlatform();
      assert.equal(result, "pc-windows-msvc");
    });

    it("should return undefined for unsupported platform", () => {
      restore = mockPlatform("unsupported");
      const result = getPlatform();
      assert.equal(result, undefined);
    });
  });
});
