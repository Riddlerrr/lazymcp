import { defineConfig } from 'vocs'

export default defineConfig({
  title: 'LazyMCP',
  description: 'A general-purpose MCP (Model Context Protocol) server written in Go that provides calculator, network, and weather tools with real-time data access.',
  sidebar: [
    {
      text: 'Getting Started',
      link: '/getting-started',
    },
    {
      text: 'Configuration',
      link: '/configuration',
    },
    {
      text: 'Tools',
      items: [
        {
          text: 'Overview',
          link: '/tools',
        },
        {
          text: 'Weather',
          link: '/tools/weather',
        },
        {
          text: 'Network',
          link: '/tools/network',
        },
        {
          text: 'Calculator',
          link: '/tools/calculator',
        },
      ],
    },
    {
      text: 'Development',
      link: '/development',
    },
  ],
})
