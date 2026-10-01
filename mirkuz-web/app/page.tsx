import React from "react";
import { Navbar } from "@/components/Navbar";
import { HeroSection } from "@/components/HeroSection";
import { StatsBar } from "@/components/StatsBar";
import { StreamCurriculum } from "@/components/StreamCurriculum";
import { FeaturesBento } from "@/components/FeaturesBento";
import { InteractiveExamDemo } from "@/components/InteractiveExamDemo";
import { OfflineShowcase } from "@/components/OfflineShowcase";
import { PricingSection } from "@/components/PricingSection";
import { Testimonials } from "@/components/Testimonials";
import { FaqSection } from "@/components/FaqSection";
import { CtaBanner } from "@/components/CtaBanner";
import { Footer } from "@/components/Footer";

export default function LandingPage() {
  return (
    <div className="min-h-screen bg-[#050505] text-[#f4f4f5] selection:bg-[#10b981]/30 selection:text-white relative">
      {/* Sticky Top Navigation */}
      <Navbar />

      {/* Main Sections */}
      <main>
        {/* 1. Hero Section */}
        <HeroSection />

        {/* 2. Trust & Metrics Bar */}
        <StatsBar />

        {/* 3. Interactive Stream & Subject Curriculum */}
        <StreamCurriculum />

        {/* 4. Core Features Bento Grid */}
        <FeaturesBento />

        {/* 5. Live Interactive Mock Exam Question Micro-Experience */}
        <InteractiveExamDemo />

        {/* 6. 100% Offline Study Architecture Showcase */}
        <OfflineShowcase />

        {/* 7. ETB Pricing & Local Telebirr/Chapa Gateways */}
        <PricingSection />

        {/* 8. Student Success Stories & Placements */}
        <Testimonials />

        {/* 9. Frequently Asked Questions */}
        <FaqSection />

        {/* 10. Bottom High-Conversion Call To Action */}
        <CtaBanner />
      </main>

      {/* Footer */}
      <Footer />
    </div>
  );
}
