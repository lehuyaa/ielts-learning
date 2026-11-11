"""Checks that every A1 bundle is imported and every lesson audio file exists.

For each bundle it calls the import *preview* (nothing is written) and reports
what an import would still create: anything to CREATE means that file (or part
of it) is not imported yet. It also lists warnings (e.g. a scenario without its
model dialogue, an unknown vocabulary slug) and missing .m4a files.

Usage (from content/a1, backend running):
    python3 check_import.py --email you@example.com
    python3 check_import.py --email you@example.com --api http://localhost:8081/api/v1
The password is asked for (or set CHECK_PASSWORD). The account must be an admin.
"""
import argparse
import getpass
import json
import os
import sys
import urllib.error
import urllib.request
import uuid

# Import order (docs/admin-guide.md): grammar and sounds first, scenarios
# before and again after the weeks (to link each scenario to its dialogue).
BUNDLES = ["a1-grammar.json", "a1-sounds.json", "a1-scenarios.json",
           "week-1.json", "week-2.json", "week-3.json", "week-4.json"]
AUDIO_ROOT = os.path.join("..", "..", "frontend", "public")


def request(url, data=None, headers=None):
    req = urllib.request.Request(url, data=data, headers=headers or {}, method="POST" if data is not None else "GET")
    try:
        with urllib.request.urlopen(req, timeout=120) as response:
            return json.load(response)
    except urllib.error.HTTPError as error:
        body = error.read().decode("utf-8", "replace")
        sys.exit(f"{url}: HTTP {error.code} {body[:300]}")
    except urllib.error.URLError as error:
        sys.exit(f"{url}: {error.reason} — is the backend running?")


def login(api, email, password):
    body = json.dumps({"email": email, "password": password}).encode()
    data = request(f"{api}/auth/login", body, {"Content-Type": "application/json"})["data"]
    return data["accessToken"]


def preview(api, token, path):
    boundary = uuid.uuid4().hex
    with open(path, "rb") as handle:
        content = handle.read()
    body = (f"--{boundary}\r\nContent-Disposition: form-data; name=\"file\"; filename=\"{os.path.basename(path)}\"\r\n"
            f"Content-Type: application/json\r\n\r\n").encode() + content + f"\r\n--{boundary}--\r\n".encode()
    headers = {"Authorization": f"Bearer {token}", "Content-Type": f"multipart/form-data; boundary={boundary}"}
    return request(f"{api}/admin/content/import/preview", body, headers)["data"]


def missing_audio():
    missing = []
    for name in sorted(os.listdir(".")):
        if not (name.startswith("audio_manifest_") and name.endswith(".tsv")):
            continue
        with open(name, encoding="utf-8") as handle:
            for line in handle:
                path = line.split("\t", 1)[0].strip()
                if path and not os.path.exists(os.path.join(AUDIO_ROOT, "audio", "a1", f"{path}.m4a")):
                    missing.append(f"{name}: {path}.m4a")
    return missing


def main():
    parser = argparse.ArgumentParser(description=__doc__, formatter_class=argparse.RawDescriptionHelpFormatter)
    parser.add_argument("--api", default="http://localhost:8081/api/v1")
    parser.add_argument("--email", required=True, help="admin account")
    args = parser.parse_args()
    password = os.environ.get("CHECK_PASSWORD") or getpass.getpass(f"Password for {args.email}: ")
    token = login(args.api.rstrip("/"), args.email, password)

    problems = 0
    for name in BUNDLES:
        if not os.path.exists(name):
            print(f"✗ {name}: file not found (run its build script)")
            problems += 1
            continue
        result = preview(args.api.rstrip("/"), token, name)
        to_create = [item for item in result.get("items", []) if item.get("action") == "CREATE"]
        invalid = [item for item in result.get("items", []) if item.get("errors")]
        warnings = result.get("warnings") or []
        errors = result.get("bundleErrors") or []
        if not to_create and not invalid and not warnings and not errors:
            print(f"✓ {name}: imported")
            continue
        problems += 1
        print(f"✗ {name}:")
        if to_create:
            kinds = {}
            for item in to_create:
                kinds[item.get("kind", "?")] = kinds.get(item.get("kind", "?"), 0) + 1
            summary = ", ".join(f"{count} {kind.lower()}" for kind, count in sorted(kinds.items()))
            print(f"    not imported yet: {summary}  → import this file")
            for item in to_create[:5]:
                print(f"      - {item.get('path')}")
        for message in errors:
            print(f"    error: {message}")
        for item in invalid[:5]:
            print(f"    invalid: {item.get('path')}: {'; '.join(item['errors'])}")
        for message in warnings[:10]:
            print(f"    warning: {message}")
        if any("model dialogue" in message for message in warnings):
            print("    → import the week bundles, then a1-scenarios.json again")

    audio = missing_audio()
    if audio:
        problems += 1
        print(f"✗ audio: {len(audio)} file(s) missing → run ./generate_lesson_audio.sh")
        for path in audio[:10]:
            print(f"    - {path}")
    else:
        print("✓ audio: every file in audio_manifest_*.tsv exists")

    print("\nAll good." if problems == 0 else f"\n{problems} item(s) to fix.")
    sys.exit(1 if problems else 0)


if __name__ == "__main__":
    main()
