"use client";

import React, { useState } from "react";
import Link from "next/link";
import { Send, CheckCircle2, Loader2, AlertCircle } from "lucide-react";

export function Footer() {
  const [email, setEmail] = useState("");
  const [loading, setLoading] = useState(false);
  const [statusMessage, setStatusMessage] = useState<{
    type: "success" | "error";
    text: string;
  } | null>(null);

  const handleSubscribe = async (e: React.FormEvent) => {
    e.preventDefault();
    if (!email) return;

    setLoading(true);
    setStatusMessage(null);

    try {
      const res = await fetch("/api/newsletter", {
        method: "POST",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify({ email }),
      });

      const data = await res.json();

      if (!res.ok || !data.success) {
        throw new Error(data.error || "Failed to subscribe.");
      }

      setStatusMessage({
        type: "success",
        text: "Welcome email sent! Thank you for subscribing.",
      });
      setEmail("");
    } catch (err: any) {
      setStatusMessage({
        type: "error",
        text: err.message || "Subscription failed. Please try again.",
      });
    } finally {
      setLoading(false);
    }
  };

  return (
    <footer className="bg-[#101730] text-white pt-20 pb-14 text-sm border-t border-[#1E2B58]">
      <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
        
        <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-12 gap-12 lg:gap-8 pb-14 border-b border-[#1E2B58]">
          
          {/* Left Column: Brand & Copyright & Socials (Col span 5) */}
          <div className="lg:col-span-5 flex flex-col justify-between space-y-8">
            
            {/* Logo */}
            <div className="space-y-4">
              <Link href="/" className="flex items-center gap-2.5 group">
                <div className="w-8 h-8 rounded-lg bg-[#00B894] flex items-center justify-center text-white shadow-sm">
                  <svg className="w-5 h-5 fill-current" viewBox="0 0 24 24">
                    <path d="M12 2L2 7l10 5 10-5-10-5zm0 9.2L4.5 7.4 12 3.6l7.5 3.8L12 11.2zm0 2.8L2 9v6l10 5 10-5V9l-10 5zm0 2.2l-7.5-3.8V11l7.5 3.8 7.5-3.8v2.6L12 16.2z"/>
                  </svg>
                </div>
                <div className="flex items-baseline gap-1.5">
                  <span className="text-2xl font-black tracking-tight text-white">
                    Mirkuz
                  </span>
                  <span className="text-xs font-bold text-[#00B894] bg-[#1E2B58] px-2.5 py-0.5 rounded-full border border-[#00B894]/30">
                    ምርኩዝ
                  </span>
                </div>
              </Link>

              <div className="space-y-1.5 text-xs text-[#94A3B8]">
                <p>Copyright © 2026 Mirkuz (ምርኩዝ) Education Platform.</p>
                <p>All rights reserved. Dedicated to Ethiopian Grade 12 students.</p>
              </div>
            </div>

            {/* Social Icons in dark rounded circles */}
            <div className="flex items-center gap-3">
              {/* Telegram Channel */}
              <a
                href="https://t.me/mirkuz_exam"
                target="_blank"
                rel="noopener noreferrer"
                className="w-10 h-10 rounded-full bg-[#1E2B58] hover:bg-[#00B894] flex items-center justify-center text-white transition-all hover:scale-105"
                aria-label="Telegram"
                title="Telegram Channel"
              >
                <Send className="w-4 h-4" />
              </a>

              {/* Instagram */}
              <a
                href="https://instagram.com"
                target="_blank"
                rel="noopener noreferrer"
                className="w-10 h-10 rounded-full bg-[#1E2B58] hover:bg-[#00B894] flex items-center justify-center text-white transition-all hover:scale-105"
                aria-label="Instagram"
              >
                <svg className="w-4 h-4 fill-current" viewBox="0 0 24 24">
                  <path d="M12 2.163c3.204 0 3.584.012 4.85.07 3.252.148 4.771 1.691 4.919 4.919.058 1.265.069 1.645.069 4.849 0 3.205-.012 3.584-.069 4.849-.149 3.225-1.664 4.771-4.919 4.919-1.266.058-1.644.07-4.85.07-3.204 0-3.584-.012-4.849-.07-3.26-.149-4.771-1.699-4.919-4.92-.058-1.265-.07-1.644-.07-4.849 0-3.204.013-3.583.07-4.849.149-3.227 1.664-4.771 4.919-4.919 1.266-.057 1.645-.069 4.849-.069zm0-2.163c-3.259 0-3.667.014-4.947.072-4.358.2-6.78 2.618-6.98 6.98-.059 1.281-.073 1.689-.073 4.948 0 3.259.014 3.668.072 4.948.2 4.358 2.618 6.78 6.98 6.98 1.281.058 1.689.072 4.948.072 3.259 0 3.668-.014 4.948-.072 4.354-.2 6.782-2.618 6.979-6.98.059-1.28.073-1.689.073-4.948 0-3.259-.014-3.667-.072-4.947-.196-4.354-2.617-6.78-6.979-6.98-1.281-.059-1.69-.073-4.949-.073zm0 5.838c-3.403 0-6.162 2.759-6.162 6.162s2.759 6.163 6.162 6.163 6.162-2.759 6.162-6.163c0-3.403-2.759-6.162-6.162-6.162zm0 10.162c-2.209 0-4-1.79-4-4 0-2.209 1.791-4 4-4s4 1.791 4 4c0 2.21-1.791 4-4 4zm6.406-11.845c-.796 0-1.441.645-1.441 1.44s.645 1.44 1.441 1.44c.795 0 1.439-.645 1.439-1.44s-.644-1.44-1.439-1.44z"/>
                </svg>
              </a>

              {/* YouTube */}
              <a
                href="https://youtube.com"
                target="_blank"
                rel="noopener noreferrer"
                className="w-10 h-10 rounded-full bg-[#1E2B58] hover:bg-[#00B894] flex items-center justify-center text-white transition-all hover:scale-105"
                aria-label="YouTube"
              >
                <svg className="w-4 h-4 fill-current" viewBox="0 0 24 24">
                  <path d="M23.498 6.186a3.016 3.016 0 0 0-2.122-2.136C19.505 3.545 12 3.545 12 3.545s-7.505 0-9.377.505A3.017 3.017 0 0 0 .502 6.186C0 8.07 0 12 0 12s0 3.93.502 5.814a3.016 3.016 0 0 0 2.122 2.136c1.871.505 9.376.505 9.376.505s7.505 0 9.377-.505a3.015 3.015 0 0 0 2.122-2.136C24 15.93 24 12 24 12s0-3.93-.502-5.814zM9.545 15.568V8.432L15.818 12l-6.273 3.568z"/>
                </svg>
              </a>
            </div>

          </div>

          {/* Column 2: Navigation Links (Col span 2) */}
          <div className="lg:col-span-2 space-y-4">
            <h4 className="text-base font-bold text-white tracking-tight">
              Platform
            </h4>
            <ul className="space-y-2.5 text-xs text-[#94A3B8]">
              <li><a href="#hero" className="hover:text-[#00B894] transition-colors">Overview</a></li>
              <li><a href="#curriculum" className="hover:text-[#00B894] transition-colors">Curriculum Streams</a></li>
              <li><a href="#features" className="hover:text-[#00B894] transition-colors">Video Lessons</a></li>
              <li><a href="#mock-exams" className="hover:text-[#00B894] transition-colors">Mock Exam Engine</a></li>
              <li><a href="#study-guides" className="hover:text-[#00B894] transition-colors">Study Guides</a></li>
            </ul>
          </div>

          {/* Column 3: Support & Legal (Col span 2) */}
          <div className="lg:col-span-2 space-y-4">
            <h4 className="text-base font-bold text-white tracking-tight">
              Support & Legal
            </h4>
            <ul className="space-y-2.5 text-xs text-[#94A3B8]">
              <li><a href="#faq" className="hover:text-[#00B894] transition-colors">Help center</a></li>
              <li><Link href="/terms" className="hover:text-[#00B894] transition-colors">Terms of service</Link></li>
              <li><Link href="/privacy" className="hover:text-[#00B894] transition-colors">Privacy policy</Link></li>
              <li><Link href="/delete-account" className="hover:text-rose-400 text-rose-400/80 transition-colors">Delete account</Link></li>
              <li><a href="mailto:support@mirkuz.app" className="hover:text-[#00B894] transition-colors">Contact support</a></li>
            </ul>
          </div>

          {/* Column 4: Stay up to date (Col span 3) */}
          <div className="lg:col-span-3 space-y-4">
            <h4 className="text-base font-bold text-white tracking-tight">
              Stay up to date
            </h4>
            
            <form onSubmit={handleSubscribe} className="relative">
              <input
                type="email"
                value={email}
                onChange={(e) => setEmail(e.target.value)}
                placeholder="Your email address"
                required
                disabled={loading}
                className="w-full bg-[#162145] text-white placeholder-[#64748B] text-xs px-4 py-3 rounded-full border border-[#1E2B58] focus:border-[#00B894] focus:outline-hidden pr-12 transition-colors disabled:opacity-60"
              />
              <button
                type="submit"
                disabled={loading}
                aria-label="Subscribe"
                className="absolute right-1.5 top-1/2 -translate-y-1/2 w-8 h-8 rounded-full bg-[#00B894] hover:bg-[#009874] flex items-center justify-center text-white transition-colors cursor-pointer disabled:opacity-50"
              >
                {loading ? (
                  <Loader2 className="w-3.5 h-3.5 animate-spin text-white" />
                ) : (
                  <Send className="w-3.5 h-3.5" />
                )}
              </button>
            </form>

            {statusMessage && (
              <div
                className={`flex items-start gap-1.5 text-xs ${
                  statusMessage.type === "success"
                    ? "text-[#00B894]"
                    : "text-rose-400"
                } animate-in fade-in duration-200`}
              >
                {statusMessage.type === "success" ? (
                  <CheckCircle2 className="w-3.5 h-3.5 mt-0.5 shrink-0" />
                ) : (
                  <AlertCircle className="w-3.5 h-3.5 mt-0.5 shrink-0" />
                )}
                <span>{statusMessage.text}</span>
              </div>
            )}
          </div>

        </div>

        {/* Bottom sub-footer */}
        <div className="pt-8 flex flex-col sm:flex-row items-center justify-between gap-4 text-xs text-[#64748B]">
          <p>© 2026 Mirkuz (ምርኩዝ) • Domain: mirkuz.app</p>
          <div className="flex items-center gap-4">
            <span>🇪🇹 Ethiopian Grade 12 National Curriculum</span>
          </div>
        </div>

      </div>
    </footer>
  );
}
