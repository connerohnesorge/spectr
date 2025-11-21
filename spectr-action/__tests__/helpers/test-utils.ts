import path from "node:path";

/**
 * Get the absolute path to a test fixture
 * @param fixtureName - Name of the fixture directory
 * @returns Absolute path to the fixture
 */
export function getFixturePath(fixtureName: string): string {
  return path.join(__dirname, "..", "fixtures", fixtureName);
}

/**
 * Create a mock process environment
 * @param env - Environment variables to set
 * @returns Function to restore original environment
 */
export function mockEnv(env: Record<string, string>): () => void {
  const originalEnv = { ...process.env };

  Object.assign(process.env, env);

  return () => {
    process.env = originalEnv;
  };
}

/**
 * Mock process.arch for testing
 * @param arch - Architecture to mock
 * @returns Function to restore original architecture
 */
export function mockArch(arch: string): () => void {
  const originalArch = process.arch;

  Object.defineProperty(process, "arch", {
    configurable: true,
    value: arch,
  });

  return () => {
    Object.defineProperty(process, "arch", {
      configurable: true,
      value: originalArch,
    });
  };
}

/**
 * Mock process.platform for testing
 * @param platform - Platform to mock
 * @returns Function to restore original platform
 */
export function mockPlatform(platform: string): () => void {
  const originalPlatform = process.platform;

  Object.defineProperty(process, "platform", {
    configurable: true,
    value: platform,
  });

  return () => {
    Object.defineProperty(process, "platform", {
      configurable: true,
      value: originalPlatform,
    });
  };
}

/**
 * Wait for a specified amount of time
 * @param ms - Milliseconds to wait
 */
export function wait(ms: number): Promise<void> {
  return new Promise((resolve) => setTimeout(resolve, ms));
}

/**
 * Assert that a value is defined (not null or undefined)
 * @param value - Value to check
 * @param message - Error message if assertion fails
 */
export function assertDefined<T>(
  value: T | null | undefined,
  message?: string,
): asserts value is T {
  if (value === null || value === undefined) {
    throw new Error(message || "Expected value to be defined");
  }
}
