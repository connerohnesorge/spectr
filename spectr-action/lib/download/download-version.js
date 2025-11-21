"use strict";
var __createBinding = (this && this.__createBinding) || (Object.create ? (function(o, m, k, k2) {
    if (k2 === undefined) k2 = k;
    var desc = Object.getOwnPropertyDescriptor(m, k);
    if (!desc || ("get" in desc ? !m.__esModule : desc.writable || desc.configurable)) {
      desc = { enumerable: true, get: function() { return m[k]; } };
    }
    Object.defineProperty(o, k2, desc);
}) : (function(o, m, k, k2) {
    if (k2 === undefined) k2 = k;
    o[k2] = m[k];
}));
var __setModuleDefault = (this && this.__setModuleDefault) || (Object.create ? (function(o, v) {
    Object.defineProperty(o, "default", { enumerable: true, value: v });
}) : function(o, v) {
    o["default"] = v;
});
var __importStar = (this && this.__importStar) || (function () {
    var ownKeys = function(o) {
        ownKeys = Object.getOwnPropertyNames || function (o) {
            var ar = [];
            for (var k in o) if (Object.prototype.hasOwnProperty.call(o, k)) ar[ar.length] = k;
            return ar;
        };
        return ownKeys(o);
    };
    return function (mod) {
        if (mod && mod.__esModule) return mod;
        var result = {};
        if (mod != null) for (var k = ownKeys(mod), i = 0; i < k.length; i++) if (k[i] !== "default") __createBinding(result, mod, k[i]);
        __setModuleDefault(result, mod);
        return result;
    };
})();
Object.defineProperty(exports, "__esModule", { value: true });
exports.tryGetFromToolCache = tryGetFromToolCache;
exports.downloadVersion = downloadVersion;
exports.resolveVersion = resolveVersion;
const node_fs_1 = require("node:fs");
const path = __importStar(require("node:path"));
const core = __importStar(require("@actions/core"));
const tc = __importStar(require("@actions/tool-cache"));
const core_1 = require("@octokit/core");
const plugin_paginate_rest_1 = require("@octokit/plugin-paginate-rest");
const plugin_rest_endpoint_methods_1 = require("@octokit/plugin-rest-endpoint-methods");
const constants_1 = require("../utils/constants");
const PaginatingOctokit = core_1.Octokit.plugin(plugin_paginate_rest_1.paginateRest, plugin_rest_endpoint_methods_1.restEndpointMethods);
function tryGetFromToolCache(arch, version) {
    core.debug(`Trying to get spectr from tool cache for ${version}...`);
    const cachedVersions = tc.findAllVersions(constants_1.TOOL_CACHE_NAME, arch);
    core.debug(`Cached versions: ${cachedVersions}`);
    let resolvedVersion = tc.evaluateVersions(cachedVersions, version);
    if (resolvedVersion === "") {
        resolvedVersion = version;
    }
    const installedPath = tc.find(constants_1.TOOL_CACHE_NAME, resolvedVersion, arch);
    return { installedPath, version: resolvedVersion };
}
async function downloadVersion(platform, arch, version, githubToken) {
    const artifact = `spectr-${arch}-${platform}`;
    let extension = ".tar.gz";
    if (platform === "pc-windows-msvc") {
        extension = ".zip";
    }
    const downloadUrl = constructDownloadUrl(version, platform, arch);
    core.debug(`Downloading spectr from "${downloadUrl}" ...`);
    const downloadPath = await tc.downloadTool(downloadUrl, undefined, githubToken);
    core.debug(`Downloaded spectr to "${downloadPath}"`);
    const extractedDir = await extractDownloadedArtifact(version, downloadPath, extension, platform, artifact);
    const cachedToolDir = await tc.cacheDir(extractedDir, constants_1.TOOL_CACHE_NAME, version, arch);
    return { cachedToolDir, version: version };
}
function constructDownloadUrl(version, platform, arch) {
    // Spectr uses simple artifact naming: spectr-{arch}-{platform}.{ext}
    const artifact = `spectr-${arch}-${platform}`;
    let extension = ".tar.gz";
    if (platform === "pc-windows-msvc") {
        extension = ".zip";
    }
    // Spectr releases use the version tag directly
    return `https://github.com/${constants_1.OWNER}/${constants_1.REPO}/releases/download/${version}/${artifact}${extension}`;
}
async function extractDownloadedArtifact(_version, downloadPath, extension, platform, _artifact) {
    let spectrDir;
    if (platform === "pc-windows-msvc") {
        const fullPathWithExtension = `${downloadPath}${extension}`;
        await node_fs_1.promises.copyFile(downloadPath, fullPathWithExtension);
        spectrDir = await tc.extractZip(fullPathWithExtension);
        // On windows extracting the zip does not create an intermediate directory
    }
    else {
        spectrDir = await tc.extractTar(downloadPath);
        // Check if an intermediate directory was created
        const files = await node_fs_1.promises.readdir(spectrDir);
        if (files.length === 1) {
            const potentialDir = path.join(spectrDir, files[0]);
            const stat = await node_fs_1.promises.stat(potentialDir);
            if (stat.isDirectory()) {
                spectrDir = potentialDir;
            }
        }
    }
    const files = await node_fs_1.promises.readdir(spectrDir);
    core.debug(`Contents of ${spectrDir}: ${files.join(", ")}`);
    return spectrDir;
}
async function resolveVersion(versionInput, githubToken) {
    core.debug(`Resolving ${versionInput}...`);
    const version = versionInput === "latest"
        ? await getLatestVersion(githubToken)
        : versionInput;
    if (tc.isExplicitVersion(version)) {
        core.debug(`Version ${version} is an explicit version.`);
        return version;
    }
    const availableVersions = await getAvailableVersions(githubToken);
    const resolvedVersion = maxSatisfying(availableVersions, version);
    if (resolvedVersion === undefined) {
        throw new Error(`No version found for ${version}`);
    }
    core.debug(`Resolved version: ${resolvedVersion}`);
    return resolvedVersion;
}
async function getAvailableVersions(githubToken) {
    try {
        const octokit = new PaginatingOctokit({
            auth: githubToken,
        });
        return await getReleaseTagNames(octokit);
    }
    catch (err) {
        if (err.message.includes("Bad credentials")) {
            core.info("No (valid) GitHub token provided. Falling back to anonymous. Requests might be rate limited.");
            const octokit = new PaginatingOctokit();
            return await getReleaseTagNames(octokit);
        }
        throw err;
    }
}
async function getReleaseTagNames(octokit) {
    const response = await octokit.paginate(octokit.rest.repos.listReleases, {
        owner: constants_1.OWNER,
        repo: constants_1.REPO,
    });
    const releaseTagNames = response.map((release) => release.tag_name);
    if (releaseTagNames.length === 0) {
        throw Error("Github API request failed while getting releases. Check the GitHub status page for outages. Try again later.");
    }
    return response.map((release) => release.tag_name);
}
async function getLatestVersion(githubToken) {
    const octokit = new PaginatingOctokit({
        auth: githubToken,
    });
    let latestRelease;
    try {
        latestRelease = await getLatestRelease(octokit);
    }
    catch (err) {
        if (err.message.includes("Bad credentials")) {
            core.info("No (valid) GitHub token provided. Falling back to anonymous. Requests might be rate limited.");
            const octokit = new PaginatingOctokit();
            latestRelease = await getLatestRelease(octokit);
        }
        else {
            core.error("Github API request failed while getting latest release. Check the GitHub status page for outages. Try again later.");
            throw err;
        }
    }
    if (!latestRelease) {
        throw new Error("Could not determine latest release.");
    }
    return latestRelease.tag_name;
}
async function getLatestRelease(octokit) {
    const { data: latestRelease } = await octokit.rest.repos.getLatestRelease({
        owner: constants_1.OWNER,
        repo: constants_1.REPO,
    });
    return latestRelease;
}
function maxSatisfying(versions, version) {
    const maxSemver = tc.evaluateVersions(versions, version);
    if (maxSemver !== "") {
        core.debug(`Found a version that satisfies the semver range: ${maxSemver}`);
        return maxSemver;
    }
    return undefined;
}
