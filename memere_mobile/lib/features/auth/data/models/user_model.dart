import '../../domain/entities/user_entity.dart';

class UserModel extends UserEntity {
  const UserModel({
    required super.id,
    required super.email,
    required super.firstName,
    required super.lastName,
    required super.role,
    super.phone,
    super.avatarUrl,
    super.isEmailVerified,
  });

  factory UserModel.fromJson(Map<String, dynamic> json) {
    final phone = json['phone'];
    final avatarUrl = json['avatar_url'];

    return UserModel(
      id: json['id'] as String,
      email: json['email'] as String,
      firstName: json['first_name'] as String? ?? '',
      lastName: json['last_name'] as String? ?? '',
      role: _parseRole(json['role'] as String? ?? 'student'),
      phone: phone is String ? phone : null,
      avatarUrl: avatarUrl is String ? avatarUrl : null,
      isEmailVerified: json['is_email_verified'] as bool? ?? false,
    );
  }

  Map<String, dynamic> toJson() => {
        'id': id,
        'email': email,
        'first_name': firstName,
        'last_name': lastName,
        'role': role.name,
        'phone': phone,
        'avatar_url': avatarUrl,
        'is_email_verified': isEmailVerified,
      };

  static UserRole _parseRole(String role) {
    switch (role) {
      case 'teacher':
        return UserRole.teacher;
      case 'admin':
        return UserRole.admin;
      default:
        return UserRole.student;
    }
  }
}
