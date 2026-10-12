class PostTag {
  const PostTag({required this.id, required this.name});

  final int id;
  final String name;

  String get href => '/tags/$id';

  factory PostTag.fromMap(Map<String, dynamic> value) {
    return PostTag(
      id: value['id'] as int,
      name: value['name'] as String,
    );
  }

  Map<String, dynamic> toMap() => {'id': id, 'name': name};
}
