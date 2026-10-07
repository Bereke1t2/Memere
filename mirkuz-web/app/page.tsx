import React from "react";
import { Navbar } from "@/components/Navbar";
import { HeroSection } from "@/components/HeroSection";
import { ClientsBar } from "@/components/ClientsBar";
import { CommunitySection } from "@/components/CommunitySection";
import { FeatureSplitOne } from "@/components/FeatureSplitOne";
import { StatsCounter } from "@/components/StatsCounter";
import { FeatureSplitTwo } from "@/components/FeatureSplitTwo";
import { AppShowcase } from "@/components/AppShowcase";
import { TestimonialSpotlight } from "@/components/TestimonialSpotlight";
import { BlogMarketingSection } from "@/components/BlogMarketingSection";
import { InteractiveExamDemo } from "@/components/InteractiveExamDemo";
import { StreamCurriculum } from "@/components/StreamCurriculum";
import { PricingSection } from "@/components/PricingSection";
import { FaqSection } from "@/components/FaqSection";
import { CtaBanner } from "@/components/CtaBanner";
import { Footer } from "@/components/Footer";

export default function LandingPage() {
  return (
    <div className="min-h-screen bg-white text-[#4D4D4D] selection:bg-[#4CAF4F]/20 selection:text-[#263238] relative">
      {/* 1. Header / Navbar */}
      <Navbar />

      {/* Main Page Flow matching Figma Minimal Landing Page */}
      <main>
        {/* 2. Hero Section: Grade 12 National University Entrance Exam */}
        <HeroSection />

        {/* 3. Partner Universities & Payment Gateways */}
        <ClientsBar />

        {/* 4. Stream Streams: Natural, Social, Offline Mode */}
        <CommunitySection />

        {/* 5. Feature Split 1: High-Yield Video Lessons by Senior Educators */}
        <FeatureSplitOne />

        {/* 6. Stats & Achievements: Empowering Students to Reach Top Universities */}
        <StatsCounter />

        {/* 7. Feature Split 2: Timed Mock Entrance Simulator */}
        <FeatureSplitTwo />

        {/* 8. Real App Screenshots & Interactive Showcase */}
        <AppShowcase />

        {/* 9. Top Scorer Spotlight: Kalkidan Bekele (642/700 - AAU Medicine) */}
        <TestimonialSpotlight />

        {/* 10. National Exam Insights & Study Guides */}
        <BlogMarketingSection />

        {/* 11. Interactive Exam Simulator Micro-Experience */}
        <InteractiveExamDemo />

        {/* 12. Stream & Curriculum Breakdown */}
        <StreamCurriculum />

        {/* 13. Transparent ETB Pricing & Local Gateways */}
        <PricingSection />

        {/* 14. Frequently Asked Questions */}
        <FaqSection />

        {/* 15. Bottom High-Impact CTA: Ready to Secure Your University Placement? */}
        <CtaBanner />
      </main>

      {/* 16. Minimal Dark Footer */}
      <Footer />
    </div>
  );
}
