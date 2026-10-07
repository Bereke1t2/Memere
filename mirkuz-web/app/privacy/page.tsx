import React from "react";
import Link from "next/link";
import { GraduationCap, ShieldCheck, ArrowLeft, Lock, Trash2, Mail, CheckCircle2 } from "lucide-react";

export const metadata = {
  title: "Privacy Policy — Mirkuz (ምርኩዝ)",
  description: "Official Privacy Policy and Data Safety Disclosures for the Mirkuz Education & Grade 12 Learning Platform.",
};

export default function PrivacyPage() {
  return (
    <div className="min-h-screen bg-[#F8FAFC] text-[#64748B]">
      {/* Top Header */}
      <header className="border-b border-[#E2E8F0] bg-white/95 backdrop-blur sticky top-0 z-50">
        <div className="max-w-5xl mx-auto px-4 sm:px-6 py-4 flex items-center justify-between">
          <Link href="/" className="flex items-center gap-3">
            <div className="w-9 h-9 rounded-xl bg-[#00B894] flex items-center justify-center text-white font-bold shadow-xs">
              <GraduationCap className="w-5 h-5" />
            </div>
            <div>
              <span className="text-xl font-black tracking-tight text-[#1E2B58]">Mirkuz</span>
              <span className="ml-2 px-2 py-0.5 text-[10px] font-bold bg-[#E8F8F5] text-[#00B894] rounded-full border border-[#00B894]/20">
                ምርኩዝ
              </span>
            </div>
          </Link>
          <Link
            href="/"
            className="inline-flex items-center gap-2 text-xs font-bold text-white bg-[#1E2B58] hover:bg-[#152042] transition-colors px-4 py-2 rounded-full shadow-xs"
          >
            <ArrowLeft className="w-3.5 h-3.5" />
            Back to Home
          </Link>
        </div>
      </header>

      {/* Main Content */}
      <main className="max-w-4xl mx-auto px-4 sm:px-6 py-12 space-y-10">
        {/* Title banner */}
        <div className="border-b border-[#E2E8F0] pb-8 space-y-3">
          <div className="inline-flex items-center gap-2 px-3.5 py-1 rounded-full bg-[#E8F8F5] border border-[#00B894]/25 text-[#00B894] text-xs font-bold">
            <ShieldCheck className="w-4 h-4" />
            Google Play Data Safety Compliant
          </div>
          <h1 className="text-3xl sm:text-4xl font-extrabold text-[#1E2B58] tracking-tight">
            Privacy Policy & Data Protection
          </h1>
          <p className="text-sm text-[#94A3B8]">
            Last updated: October 2026 • Applies to Mirkuz Mobile App (Android/iOS) and Web Services
          </p>
        </div>

        {/* Section 1: Overview */}
        <section className="space-y-4">
          <h2 className="text-xl font-bold text-[#1E2B58] flex items-center gap-2">
            <span className="text-[#00B894]">1.</span> Overview & Commitment
          </h2>
          <p className="text-sm text-[#475569] leading-relaxed">
            Welcome to <strong className="text-[#1E2B58]">Mirkuz (ምርኩዝ)</strong>, an interactive learning platform designed to help Ethiopian secondary students excel in the Ethiopian Grade 12 National University Entrance Examination. We respect your privacy and are committed to protecting any personal information collected when using our mobile application and digital learning services.
          </p>
          <p className="text-sm text-[#475569] leading-relaxed">
            This Privacy Policy describes what data we collect, how we handle and protect it, and how you can exercise your rights, including requesting complete account and data deletion.
          </p>
        </section>

        {/* Section 2: Data We Collect */}
        <section className="space-y-4">
          <h2 className="text-xl font-bold text-[#1E2B58] flex items-center gap-2">
            <span className="text-[#00B894]">2.</span> Data We Collect & How It Is Used
          </h2>
          <div className="grid sm:grid-cols-2 gap-4 pt-2">
            <div className="p-5 rounded-2xl bg-white border border-[#E2E8F0] shadow-xs space-y-2">
              <h3 className="text-sm font-bold text-[#1E2B58] flex items-center gap-2">
                <CheckCircle2 className="w-4 h-4 text-[#00B894]" />
                Account Credentials
              </h3>
              <p className="text-xs text-[#64748B] leading-relaxed">
                First name, last name, and email address (and optional phone number) required to register your student profile, verify your identity, and safeguard your course enrollments.
              </p>
            </div>

            <div className="p-5 rounded-2xl bg-white border border-[#E2E8F0] shadow-xs space-y-2">
              <h3 className="text-sm font-bold text-[#1E2B58] flex items-center gap-2">
                <CheckCircle2 className="w-4 h-4 text-[#00B894]" />
                Academic & Exam Analytics
              </h3>
              <p className="text-xs text-[#64748B] leading-relaxed">
                Quiz and mock entrance exam answers, completion timestamps, subject percentage scores, and lesson progress to compute personalized study metrics and performance reports.
              </p>
            </div>

            <div className="p-5 rounded-2xl bg-white border border-[#E2E8F0] shadow-xs space-y-2">
              <h3 className="text-sm font-bold text-[#1E2B58] flex items-center gap-2">
                <CheckCircle2 className="w-4 h-4 text-[#00B894]" />
                Device & Session Security
              </h3>
              <p className="text-xs text-[#64748B] leading-relaxed">
                Unique client device identifiers to ensure single-session login integrity during national mock entrance exams and avoid unauthorized multi-device account sharing.
              </p>
            </div>

            <div className="p-5 rounded-2xl bg-white border border-[#E2E8F0] shadow-xs space-y-2">
              <h3 className="text-sm font-bold text-[#1E2B58] flex items-center gap-2">
                <CheckCircle2 className="w-4 h-4 text-[#00B894]" />
                Offline Storage & Caching
              </h3>
              <p className="text-xs text-[#64748B] leading-relaxed">
                Local device encrypted storage for downloaded course notes, study PDFs, video cache, and pending offline exam submissions synced when internet is restored.
              </p>
            </div>
          </div>
        </section>

        {/* Section 3: Data Security */}
        <section className="space-y-4">
          <h2 className="text-xl font-bold text-[#1E2B58] flex items-center gap-2">
            <span className="text-[#00B894]">3.</span> Data Encryption & Security Standards
          </h2>
          <div className="p-6 rounded-2xl bg-white border border-[#E2E8F0] shadow-xs space-y-3">
            <div className="flex items-center gap-3">
              <Lock className="w-5 h-5 text-[#00B894]" />
              <span className="text-sm font-bold text-[#1E2B58]">Encryption in Transit & at Rest</span>
            </div>
            <p className="text-xs text-[#64748B] leading-relaxed">
              All data transmitted between the Mirkuz mobile application and our backend servers is strictly protected with modern TLS/HTTPS encryption. Sensitive auth tokens and locally stored answer keys are protected with hardware-backed secure storage.
            </p>
            <p className="text-xs text-[#64748B] leading-relaxed">
              We never sell, rent, or trade student data with third-party advertisers or data brokers.
            </p>
          </div>
        </section>

        {/* Section 4: Third-Party Services */}
        <section className="space-y-4">
          <h2 className="text-xl font-bold text-[#1E2B58] flex items-center gap-2">
            <span className="text-[#00B894]">4.</span> Third-Party Service Providers
          </h2>
          <p className="text-sm text-[#475569] leading-relaxed">
            We utilize trusted infrastructure providers solely to deliver platform functionality:
          </p>
          <ul className="list-disc list-inside text-xs text-[#64748B] space-y-2 pl-2">
            <li><strong className="text-[#1E2B58]">Cloud Hosting & Database:</strong> Secure backend API hosting with PostgreSQL database infrastructure.</li>
            <li><strong className="text-[#1E2B58]">Payment Processing:</strong> Chapa / Telebirr payment gateways for processing course enrollments and subscription access in Ethiopian Birr (ETB). Credit card or Telebirr PIN data is handled directly by licensed payment gateways and never stored on Mirkuz servers.</li>
            <li><strong className="text-[#1E2B58]">Notifications:</strong> Firebase Cloud Messaging for study reminders and mock exam announcements.</li>
          </ul>
        </section>

        {/* Section 5: Children's Privacy */}
        <section className="space-y-4">
          <h2 className="text-xl font-bold text-[#1E2B58] flex items-center gap-2">
            <span className="text-[#00B894]">5.</span> Children & Student Privacy (COPPA / GDPR-K)
          </h2>
          <p className="text-sm text-[#475569] leading-relaxed">
            Mirkuz is designed for secondary school students preparing for university entrance. We collect only the minimal necessary academic data to provide tutoring, mock exams, and curriculum access. If a parent or guardian becomes aware that their child has provided personal information without permission, they may contact us for immediate deletion.
          </p>
        </section>

        {/* Section 6: Account Deletion */}
        <section className="space-y-4 p-6 rounded-2xl bg-white border border-[#E2E8F0] shadow-dribbble">
          <div className="flex items-center gap-3">
            <Trash2 className="w-6 h-6 text-rose-500" />
            <h2 className="text-xl font-bold text-[#1E2B58]">6. Your Rights & Account Deletion</h2>
          </div>
          <p className="text-sm text-[#475569] leading-relaxed">
            You have the right to access, export, or permanently delete your account and all associated personal data at any time.
          </p>
          <div className="pt-2 flex flex-col sm:flex-row gap-3">
            <Link
              href="/delete-account"
              className="inline-flex items-center justify-center gap-2 px-5 py-2.5 rounded-full bg-rose-500 text-white text-xs font-bold hover:bg-rose-600 transition-colors shadow-xs"
            >
              <Trash2 className="w-4 h-4" />
              Request Web Account Deletion
            </Link>
            <div className="text-xs text-[#94A3B8] flex items-center">
              Or tap <em>Profile → Settings → Account Session → Delete Account</em> directly inside the Mirkuz mobile app.
            </div>
          </div>
        </section>

        {/* Section 7: Contact */}
        <section className="space-y-3 border-t border-[#E2E8F0] pt-8">
          <h2 className="text-lg font-bold text-[#1E2B58] flex items-center gap-2">
            <Mail className="w-5 h-5 text-[#00B894]" />
            7. Contact Support & Inquiries
          </h2>
          <p className="text-xs text-[#64748B] leading-relaxed">
            If you have questions regarding this Privacy Policy or wish to make an inquiry regarding your data, reach out to our team at:
          </p>
          <div className="text-xs font-bold text-[#00B894]">
            Email: <a href="mailto:support@mirkuz.app" className="underline hover:text-[#009874]">support@mirkuz.app</a> • Addis Ababa, Ethiopia
          </div>
        </section>
      </main>

      {/* Footer */}
      <footer className="border-t border-[#E2E8F0] py-8 text-center text-xs text-[#94A3B8]">
        <p>© {new Date().getFullYear()} Mirkuz (ምርኩዝ) Education Platform. All rights reserved.</p>
      </footer>
    </div>
  );
}
