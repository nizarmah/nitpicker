package patch

import (
	"strings"
	"testing"
)

func TestParse(t *testing.T) {
	input := `From abc123 Mon Sep 17 00:00:00 2001
From: Author <author@example.com>
Date: Mon, 1 Jan 2024 00:00:00 +0000
Subject: [PATCH] Add feature

---
 src/app.tsx  | 15 ++++++++++++---
 src/index.ts | 45 +++++++++++++++++++++++++++++++++++++++++++++
 2 files changed, 57 insertions(+), 3 deletions(-)

diff --git a/src/app.tsx b/src/app.tsx
index 1234567..abcdefg 100644
--- a/src/app.tsx
+++ b/src/app.tsx
@@ -1,5 +1,17 @@
 import React from 'react';

-function App() {
-  return <div>Hello</div>;
-}
+function App() {
+  return (
+    <div>
+      <h1>Hello World</h1>
+      <p>Welcome to the app</p>
+      <ul>
+        <li>Item 1</li>
+        <li>Item 2</li>
+        <li>Item 3</li>
+      </ul>
+    </div>
+  );
+}

 export default App;
diff --git a/src/index.ts b/src/index.ts
new file mode 100644
index 0000000..1234567
--- /dev/null
+++ b/src/index.ts
@@ -0,0 +1,5 @@
+import App from './app';
+
+const root = document.getElementById('root');
+App(root);
+console.log('started');
`

	stats, err := Parse(strings.NewReader(input))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(stats) != 2 {
		t.Fatalf("expected 2 files, got %d", len(stats))
	}

	// src/app.tsx: 13 added, 3 deleted
	if stats[0].Path != "src/app.tsx" {
		t.Errorf("expected path src/app.tsx, got %s", stats[0].Path)
	}
	if stats[0].Added != 13 {
		t.Errorf("src/app.tsx: expected 13 added, got %d", stats[0].Added)
	}
	if stats[0].Deleted != 3 {
		t.Errorf("src/app.tsx: expected 3 deleted, got %d", stats[0].Deleted)
	}

	// src/index.ts: 5 added, 0 deleted
	if stats[1].Path != "src/index.ts" {
		t.Errorf("expected path src/index.ts, got %s", stats[1].Path)
	}
	if stats[1].Added != 5 {
		t.Errorf("src/index.ts: expected 5 added, got %d", stats[1].Added)
	}
	if stats[1].Deleted != 0 {
		t.Errorf("src/index.ts: expected 0 deleted, got %d", stats[1].Deleted)
	}
}

func TestParseSameFileMultipleCommits(t *testing.T) {
	input := `diff --git a/README.md b/README.md
--- a/README.md
+++ b/README.md
@@ -1 +1,2 @@
 # Hello
+World
diff --git a/README.md b/README.md
--- a/README.md
+++ b/README.md
@@ -1,2 +1,3 @@
 # Hello
 World
+Again
`

	stats, err := Parse(strings.NewReader(input))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(stats) != 1 {
		t.Fatalf("expected 1 file, got %d", len(stats))
	}
	if stats[0].Added != 2 {
		t.Errorf("expected 2 added, got %d", stats[0].Added)
	}
	if stats[0].Deleted != 0 {
		t.Errorf("expected 0 deleted, got %d", stats[0].Deleted)
	}
}
