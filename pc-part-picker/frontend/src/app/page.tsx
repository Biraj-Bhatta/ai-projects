"use client";

import { useState, useEffect } from "react";

type Component = {
  id: number;
  name: string;
  type: string;
  price: number;
};

export default function Home() {
  const [workload, setWorkload] = useState("gaming");
  const [resolution, setResolution] = useState("1080p");
  const [colorDepth, setColorDepth] = useState("8bit");
  const [gpuPref, setGpuPref] = useState("discrete");
  const [coolingPref, setCoolingPref] = useState("air");
  const [inCountry, setInCountry] = useState(true);
  const [budget, setBudget] = useState(1000);

  const [components, setComponents] = useState<Component[]>([]);
  const [lockedComponents, setLockedComponents] = useState<number[]>([]);
  const [buildResult, setBuildResult] = useState<{ build: Component[], total_cost: number } | null>(null);
  const [loading, setLoading] = useState(false);

  useEffect(() => {
    fetch("http://localhost:8080/api/components")
      .then((res) => res.json())
      .then((data) => setComponents(data))
      .catch((err) => console.error(err));
  }, []);

  const handleGenerateBuild = async () => {
    setLoading(true);
    try {
      const res = await fetch("http://localhost:8080/api/build", {
        method: "POST",
        headers: {
          "Content-Type": "application/json",
        },
        body: JSON.stringify({
          workload,
          in_country: inCountry,
          preferences_cooling: coolingPref,
          preferences_gpu: gpuPref,
          resolution,
          color_depth: colorDepth,
          budget: Number(budget),
          locked_components: lockedComponents,
        }),
      });
      const data = await res.json();
      setBuildResult(data);
    } catch (error) {
      console.error(error);
    }
    setLoading(false);
  };

  const toggleLock = (id: number) => {
    setLockedComponents((prev) =>
      prev.includes(id) ? prev.filter((c) => c !== id) : [...prev, id]
    );
  };

  return (
    <div className="min-h-screen p-8 font-sans">
      <h1 className="text-3xl font-bold mb-6">AI PC Part Picker</h1>

      <div className="grid grid-cols-1 md:grid-cols-2 gap-8">
        {/* Configuration Panel */}
        <div className="bg-gray-100 p-6 rounded-lg text-black">
          <h2 className="text-xl font-semibold mb-4">Configuration</h2>

          <div className="mb-4">
            <label className="block mb-1">Workload</label>
            <select className="w-full p-2 border rounded" value={workload} onChange={(e) => setWorkload(e.target.value)}>
              <option value="gaming">Gaming-heavy</option>
              <option value="editing">Editing-heavy</option>
              <option value="programming">Programming</option>
              <option value="general">General Use</option>
            </select>
          </div>

          <div className="mb-4 flex gap-4">
            <div className="flex-1">
              <label className="block mb-1">Resolution</label>
              <select className="w-full p-2 border rounded" value={resolution} onChange={(e) => setResolution(e.target.value)}>
                <option value="1080p">1080p</option>
                <option value="1440p">1440p</option>
                <option value="4k">4K</option>
              </select>
            </div>
            <div className="flex-1">
              <label className="block mb-1">Color Depth</label>
              <select className="w-full p-2 border rounded" value={colorDepth} onChange={(e) => setColorDepth(e.target.value)}>
                <option value="8bit">8-bit</option>
                <option value="10bit">10-bit</option>
              </select>
            </div>
          </div>

          <div className="mb-4 flex gap-4">
            <div className="flex-1">
              <label className="block mb-1">GPU Preference</label>
              <select className="w-full p-2 border rounded" value={gpuPref} onChange={(e) => setGpuPref(e.target.value)}>
                <option value="discrete">Discrete GPU</option>
                <option value="igpu">Integrated Graphics (iGPU)</option>
              </select>
            </div>
            <div className="flex-1">
              <label className="block mb-1">Cooling</label>
              <select className="w-full p-2 border rounded" value={coolingPref} onChange={(e) => setCoolingPref(e.target.value)}>
                <option value="air">Air Cooling</option>
                <option value="liquid">Liquid Cooling (AIO)</option>
              </select>
            </div>
          </div>

          <div className="mb-4">
            <label className="block mb-1">Budget ($)</label>
            <input type="number" className="w-full p-2 border rounded" value={budget} onChange={(e) => setBudget(Number(e.target.value))} />
          </div>

          <div className="mb-4 flex items-center">
            <input type="checkbox" id="inCountry" className="mr-2" checked={inCountry} onChange={(e) => setInCountry(e.target.checked)} />
            <label htmlFor="inCountry">In-Country Pricing (Uncheck for Cross-Border)</label>
          </div>

          <div className="mb-4">
            <h3 className="font-semibold mb-2">Lock Components</h3>
            <div className="max-h-40 overflow-y-auto border p-2 rounded bg-white">
              {components.map((comp) => (
                <div key={comp.id} className="flex items-center mb-1">
                  <input
                    type="checkbox"
                    id={`comp-${comp.id}`}
                    className="mr-2"
                    checked={lockedComponents.includes(comp.id)}
                    onChange={() => toggleLock(comp.id)}
                  />
                  <label htmlFor={`comp-${comp.id}`}>
                    [{comp.type}] {comp.name}
                  </label>
                </div>
              ))}
            </div>
          </div>

          <button
            className="w-full bg-blue-600 text-white p-3 rounded font-bold hover:bg-blue-700"
            onClick={handleGenerateBuild}
            disabled={loading}
          >
            {loading ? "Generating..." : "Generate Build"}
          </button>
        </div>

        {/* Results Panel */}
        <div className="bg-gray-100 p-6 rounded-lg text-black">
          <h2 className="text-xl font-semibold mb-4">Recommended Build</h2>
          {buildResult ? (
            <div>
              <div className="space-y-4 mb-6">
                {buildResult.build.map((part, idx) => (
                  <div key={idx} className="bg-white p-4 rounded shadow flex justify-between items-center">
                    <div>
                      <div className="text-sm text-gray-500 font-bold">{part.type}</div>
                      <div className="font-semibold">{part.name}</div>
                    </div>
                    <div className="text-green-600 font-bold">
                      ${part.price.toFixed(2)}
                    </div>
                  </div>
                ))}
              </div>
              <div className="text-right text-2xl font-bold border-t pt-4">
                Total: ${buildResult.total_cost.toFixed(2)}
              </div>
            </div>
          ) : (
            <div className="text-gray-500 italic">Configure your preferences and generate a build to see recommendations here.</div>
          )}
        </div>
      </div>
    </div>
  );
}
