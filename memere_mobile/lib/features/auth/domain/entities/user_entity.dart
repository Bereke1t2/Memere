/// Pure Dart — zero external dependencies
class UserEntity {
  const UserEntity({
    required this.id,
    required this.email,
    required this.firstName,
    required this.lastName,
    required this.role,
    this.phone,
    this.avatarUrl,
    this.isEmailVerified = false,
  });

  final String id;
  final String email;
  final String firstName;
  final String lastName;
  final UserRole role;
  final String? phone;
  final String? avatarUrl;
  final bool isEmailVerified;

  String get fullName => '$firstName $lastName';
  bool get isStudent => role == UserRole.student;
  bool get isTeacher => role == UserRole.teacher;
  bool get isAdmin   => role == UserRole.admin;
}

enum UserRole { student, teacher, admin }
