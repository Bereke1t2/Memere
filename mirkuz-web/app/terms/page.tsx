import React from "react";
import Link from "next/link";
import Image from "next/image";
import { ArrowLeft, BookOpen, FileCheck, AlertCircle } from "lucide-react";

export const metadata = {
  title: "Terms of Service — Mirkuz (ምርኩዝ)",
  description: "Terms and conditions of use for the Mirkuz Grade 12 University Entrance Learning Platform.",
};

export default function TermsPage() {
  return (
    <div className="min-h-screen bg-[#F8FAFC] text-[#64748B]">
      {/* Top Header */}
      <header className="border-b border-[#E2E8F0] bg-white/95 backdrop-blur sticky top-0 z-50">
        <div className="max-w-5xl mx-auto px-4 sm:px-6 py-4 flex items-center justify-between">
          <Link href="/" className="flex items-center gap-3">
            <Image
              src="/logo.png"
              alt="Mirkuz Logo"
              width={36}
              height={36}
              className="w-9 h-9 rounded-xl shadow-xs object-contain"
            />
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
            <BookOpen className="w-4 h-4" />
            Platform Service Agreement
          </div>
          <h1 className="text-3xl sm:text-4xl font-extrabold text-[#1E2B58] tracking-tight">
            Terms of Service
          </h1>
          <p className="text-sm text-[#94A3B8]">
            Effective Date: October 2026 • Mirkuz Education Platform
          </p>
        </div>

        {/* Section 1 */}
        <section className="space-y-3">
          <h2 className="text-lg font-bold text-[#1E2B58] flex items-center gap-2">
            <span className="text-[#00B894]">1.</span> Acceptance of Terms
          </h2>
          <p className="text-xs sm:text-sm text-[#475569] leading-relaxed">
            By creating an account or accessing Mirkuz (including mobile applications, web portals, course materials, video lectures, and mock exams), you agree to be bound by these Terms of Service. If you do not agree with these terms, you should not use the platform.
          </p>
        </section>

        {/* Section 2 */}
        <section className="space-y-3">
          <h2 className="text-lg font-bold text-[#1E2B58] flex items-center gap-2">
            <span className="text-[#00B894]">2.</span> Educational Purpose & License
          </h2>
          <p className="text-xs sm:text-sm text-[#475569] leading-relaxed">
            All curriculum videos, question banks, study guides, and mock exams are protected intellectual property created for individual academic preparation for Ethiopian Grade 12 students. You are granted a limited, personal, non-exclusive, non-transferable license to access the content.
          </p>
          <div className="p-5 rounded-2xl bg-white border border-[#E2E8F0] shadow-xs space-y-2">
            <div className="flex items-center gap-2 text-xs font-bold text-amber-600">
              <AlertCircle className="w-4 h-4" />
              Prohibited Content Usage:
            </div>
            <ul className="list-disc list-inside text-xs text-[#64748B] space-y-1 pl-2">
              <li>Recording, ripping, re-uploading, or publicly broadcasting proprietary video lessons.</li>
              <li>Scraping exam question banks or distributing answer keys commercially.</li>
              <li>Sharing or renting account credentials across multiple individuals.</li>
            </ul>
          </div>
        </section>

        {/* Section 3 */}
        <section className="space-y-3">
          <h2 className="text-lg font-bold text-[#1E2B58] flex items-center gap-2">
            <span className="text-[#00B894]">3.</span> Payments & Access Subscriptions
          </h2>
          <p className="text-xs sm:text-sm text-[#475569] leading-relaxed">
            Premium course access, full national mock entrance exams, and chapter downloads can be unlocked via supported payment gateways (e.g. Telebirr, CBE Birr, Chapa). Access remains valid for the duration of the active subscription or individual course enrollment tier.
          </p>
        </section>

        {/* Section 4 */}
        <section className="space-y-3">
          <h2 className="text-lg font-bold text-[#1E2B58] flex items-center gap-2">
            <span className="text-[#00B894]">4.</span> Account Termination & Data Deletion
          </h2>
          <p className="text-xs sm:text-sm text-[#475569] leading-relaxed">
            You may terminate your account at any time via the mobile app profile settings or our{" "}
            <Link href="/delete-account" className="text-[#00B894] font-bold hover:underline">
              Account Deletion Page
            </Link>. We reserve the right to suspend or terminate accounts that engage in harassment, unauthorized content redistribution, or fraudulent activity.
          </p>
        </section>

        {/* Section 5 */}
        <section className="space-y-3 border-t border-[#E2E8F0] pt-8">
          <h2 className="text-lg font-bold text-[#1E2B58] flex items-center gap-2">
            <FileCheck className="w-5 h-5 text-[#00B894]" />
            5. Governing Law
          </h2>
          <p className="text-xs text-[#64748B] leading-relaxed">
            These terms are governed by the laws and regulations of the Federal Democratic Republic of Ethiopia.
          </p>
          <p className="text-xs text-[#64748B]">
            Questions? Contact <a href="mailto:support@mirkuz.app" className="text-[#00B894] font-bold hover:underline">support@mirkuz.app</a>.
          </p>
        </section>
      </main>

      {/* Footer */}
      <footer className="border-t border-[#E2E8F0] py-8 text-center text-xs text-[#94A3B8]">
        <p>© {new Date().getFullYear()} Mirkuz (ምርኩዝ) Education Platform. All rights reserved.</p>
      </footer>
    </div>
  );
}
