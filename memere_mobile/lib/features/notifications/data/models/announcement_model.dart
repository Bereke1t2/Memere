class AnnouncementModel {
  const AnnouncementModel({
    required this.id,
    required this.title,
    required this.body,
    required this.type,
    this.createdAt,
  });

  final String id;
  final String title;
  final String body;
  final String type;
  final DateTime? createdAt;

  factory AnnouncementModel.fromJson(Map<String, dynamic> json) {
    return AnnouncementModel(
      id: json['id'] as String? ?? '',
      title: json['title'] as String? ?? 'Announcement',
      body: json['body'] as String? ?? '',
      type: json['type'] as String? ?? 'announcement',
      createdAt: json['created_at'] != null
          ? DateTime.tryParse(json['created_at'] as String)
          : null,
    );
  }
}
