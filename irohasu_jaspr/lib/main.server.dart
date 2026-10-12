/// The entrypoint for the **server** environment.
///
/// The [main] method will only be executed on the server during pre-rendering.
/// To run code on the client, check the `main.client.dart` file.
library;

import 'package:jaspr/dom.dart';
// Server-specific Jaspr import.
import 'package:jaspr/server.dart';

// Imports the [App] component.
import 'app.dart';

// This file is generated automatically by Jaspr, do not remove or edit.
import 'main.server.options.dart';

void main() {
  // Initializes the server environment with the generated default options.
  Jaspr.initializeApp(
    options: defaultServerOptions,
  );

  // Starts the app.
  //
  // [Document] renders the root document structure (<html>, <head> and <body>)
  // with the provided parameters and components.
  runApp(
    Document(
      title: 'Irohasu',
      styles: [
        // Special import rule to include to another css file.
        css.import(
          'https://fonts.googleapis.com/css2?family=Caveat:wght@600&family=Inter:wght@400;500;600&display=swap',
        ),
        css.import(
          'https://fonts.googleapis.com/css2?family=Material+Symbols+Outlined:opsz,wght,FILL,GRAD@20..48,400,0,0',
        ),
        // Each style rule takes a valid css selector and a set of styles.
        // Styles are defined using type-safe css bindings and can be freely chained and nested.
        css('html, body').styles(
          width: 100.percent,
          minHeight: 100.vh,
          padding: .zero,
          margin: .zero,
          fontFamily: const .list([FontFamily('Inter'), FontFamilies.sansSerif]),
        ),
      ],
      head: [
        link(rel: 'icon', type: 'image/png', href: 'assets/icon.png'),
        // Keep the web app manifest. Flutter's bootstrap script is intentionally
        // omitted because this page does not mount an embedded Flutter component.
        link(rel: 'manifest', href: 'manifest.json'),
      ],
      body: App(),
    ),
  );
}
