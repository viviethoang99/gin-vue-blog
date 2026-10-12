import '../data/model/post.dart';

class TagSummary {
  const TagSummary({required this.id, required this.name, required this.count});

  final int id;
  final String name;
  final int count;
}

enum TagSort { count, name }

List<TagSummary> collectTagSummaries(Iterable<Post> posts) {
  final counts = <String, int>{};
  final ids = <String, int>{};
  final names = <String, String>{};

  for (final post in posts) {
    final tagsInPost = <String>{};
    for (final tag in post.tags ?? const <Category>[]) {
      final name = tag.name.trim();
      if (name.isEmpty) continue;

      final key = name.toLowerCase();
      ids.putIfAbsent(key, () => tag.id);
      names.putIfAbsent(key, () => name);
      tagsInPost.add(key);
    }

    for (final key in tagsInPost) {
      counts.update(key, (count) => count + 1, ifAbsent: () => 1);
    }
  }

  final summaries = [
    for (final entry in counts.entries) TagSummary(id: ids[entry.key]!, name: names[entry.key]!, count: entry.value),
  ];
  return sortTagSummaries(summaries, TagSort.count);
}

List<TagSummary> sortTagSummaries(Iterable<TagSummary> tags, TagSort sort) {
  final sorted = tags.toList(growable: false);
  sorted.sort((a, b) {
    final byName = a.name.toLowerCase().compareTo(b.name.toLowerCase());
    if (sort == TagSort.name) return byName;

    final byCount = b.count.compareTo(a.count);
    return byCount != 0 ? byCount : byName;
  });
  return sorted;
}
