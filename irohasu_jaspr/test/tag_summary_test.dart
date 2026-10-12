import 'package:flutter_test/flutter_test.dart';
import 'package:irohasu_jaspr/data/model/post.dart';
import 'package:irohasu_jaspr/services/tag_summary.dart';

void main() {
  test('collectTagSummaries counts posts, merges casing, and sorts by count', () {
    final timestamp = DateTime.utc(2026);
    Category tag(int id, String name) => Category(id: id, createdAt: timestamp, updatedAt: timestamp, name: name);

    final summaries = collectTagSummaries([
      Post(tags: [tag(1, 'Dart'), tag(2, 'Jaspr'), tag(3, 'dart')]),
      Post(tags: [tag(4, ' dart '), tag(5, 'UI')]),
      Post(tags: [tag(6, 'UI'), tag(7, '  ')]),
    ]);

    expect(
      summaries.map((tag) => (tag.name, tag.count)),
      [('Dart', 2), ('UI', 2), ('Jaspr', 1)],
    );
  });

  test('sortTagSummaries sorts alphabetically by name', () {
    final sorted = sortTagSummaries(
      const [
        TagSummary(id: 1, name: 'Vue', count: 3),
        TagSummary(id: 2, name: 'Dart', count: 1),
        TagSummary(id: 3, name: 'jaspr', count: 2),
      ],
      TagSort.name,
    );

    expect(sorted.map((tag) => tag.name), ['Dart', 'jaspr', 'Vue']);
  });
}
