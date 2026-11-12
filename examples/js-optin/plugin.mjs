import { serve } from "../../js/index.js";

serve({
  info: () => ({
    id: "com.example.jsoptin",
    name: "JavaScript Example",
    version: "0.1.0",
    api_version: 1,
    targets: ["course.new", "course.edit"],
    events: ["course.create", "course.update"],
  }),

  contribute: () => ({
    sections: [{
      title: "JavaScript Options",
      fields: [{
        key: "js_flag",
        label: "Enable the JavaScript feature",
        type: "checkbox",
      }],
    }],
  }),

  handle: (event) => {
    console.error(`js extension: ${event.phase} ${event.name}`);
    return {};
  },
});
