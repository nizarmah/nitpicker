package patch

import (
	"strings"
	"testing"
)

func TestParse(t *testing.T) {
	input := `diff --git a/src/app.tsx b/src/app.tsx
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

func TestParseMultipleHunks(t *testing.T) {
	// Simulates the squashed diff from nizarmah/mm.quest/pull/4.diff
	input := `diff --git a/js/main.js b/js/main.js
index f42e856..254ec71 100644
--- a/js/main.js
+++ b/js/main.js
@@ -180,13 +180,14 @@ const goToLoader = (screen) => {
   const quote = document.createElement("div")
   quote.className = "quote"

-  const vanillaSkies = [
-    "it's the little things.",
-    "there's nothing bigger...",
-    "is there?"
+  const ferrisBueller = [
+    "Life moves pretty fast.",
+    "If you don't stop and",
+    "look around once in a while,",
+    "you could miss it."
   ]

-  vanillaSkies.forEach((line) => {
+  ferrisBueller.forEach((line) => {
     const span = document.createElement("span")
     span.textContent = line
     quote.appendChild(span)
`

	stats, err := Parse(strings.NewReader(input))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(stats) != 1 {
		t.Fatalf("expected 1 file, got %d", len(stats))
	}
	if stats[0].Path != "js/main.js" {
		t.Errorf("expected path js/main.js, got %s", stats[0].Path)
	}
	if stats[0].Added != 6 {
		t.Errorf("expected 6 added, got %d", stats[0].Added)
	}
	if stats[0].Deleted != 5 {
		t.Errorf("expected 5 deleted, got %d", stats[0].Deleted)
	}
}
