/// Pure Dart — zero external dependencies
class UserEntity {
  const UserEntity({
    required this.id,
    required this.email,
    required this.firstName,
    required this.lastName,
    required this.role,
    this.approvalStatus = UserApprovalStatus.approved,
    this.phone,
    this.avatarUrl,
    this.isEmailVerified = false,
  });

  final String id;
  final String email;
  final String firstName;
  final String lastName;
  final UserRole role;
  final UserApprovalStatus approvalStatus;
  final String? phone;
  final String? avatarUrl;
  final bool isEmailVerified;

  String get fullName => '$firstName $lastName';
  bool get isStudent => role == UserRole.student;
  bool get isTeacher => role == UserRole.teacher;
  bool get isAdmin   => role == UserRole.admin;
  bool get isApproved => approvalStatus == UserApprovalStatus.approved;
  bool get isPendingApproval => approvalStatus == UserApprovalStatus.pending;
  bool get isRejected => approvalStatus == UserApprovalStatus.rejected;
}

enum UserRole { student, teacher, admin }

enum UserApprovalStatus { pending, approved, rejected }
