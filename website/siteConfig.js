/**
 * Copyright (c) 2017-present, Facebook, Inc.
 *
 * This source code is licensed under the MIT license found in the
 * LICENSE file in the root directory of this source tree.
 */

// See https://docusaurus.io/docs/site-config for all the possible
// site configuration options.


const {Plugin: Embed} = require('remarkable-embed');

// Our custom remarkable plugin factory.
const createVariableInjectionPlugin = () => {
  // `let` binding used to initialize the `Embed` plugin only once for efficiency.
  // See `if` statement below.
  let initializedPlugin;

  const embed = new Embed();
  embed.register({
    // Call the render method to process the corresponding variable with
    // the passed Remarkable instance.
    // -> the Markdown markup in the variable will be converted to HTML.
    injectImage: path => `<img src="https://kdarkroom.herokuapp.com/${path}" alt="image"/>`,
  });

  return (md, options) => {
    if (!initializedPlugin) {
      initializedPlugin = {
        render: md.render.bind(md),
        hook: embed.hook(md, options)
      };
    }

    return initializedPlugin.hook;
  };
};

// List of projects/orgs using your project for the users page.
const users = [
  {
    caption: 'Gojek',
    // You will need to prepend the image path with your baseUrl
    // if it is not '/', like: '/test-site/img/image.jpg'.
    image: './img/gojek-logo.svg',
    infoLink: 'https://www.gojek.io',
    pinned: true,
  },
];

const siteConfig = {
  title: 'Darkroom', // Title for your website.
  tagline: 'Yet Another Image Proxy',
  githubUrl: 'https://github.com/gojek/darkroom', // Your website URL
  baseUrl: '/darkroom/', // Base URL for your project */
  // For github.io type URLs, you would set the url and baseUrl like:
  url: 'https://www.gojek.io',
 
  // Used for publishing and more
  projectName: 'darkroom',
  organizationName: 'gojek',
  // For top-level user or org sites, the organization is still the same.
  // e.g., for the https://JoelMarcey.github.io site, it would be set like...
  //   organizationName: 'JoelMarcey'

  // For no header links in the top nav bar -> headerLinks: [],
  markdownPlugins: [
    createVariableInjectionPlugin()
  ],
  headerLinks: [
    {doc: 'getting-started', label: 'Docs'},
    {blog: true, label: 'Updates'}
  ],

  // If you have users set above, you add it here:
  users,

  /* path to images for header/footer */
  headerIcon: 'img/darkroom-logo.png',
  footerIcon: 'img/darkroom-logo.png',
  favicon: 'img/darkroom-logo.png',

  /* Colors for website */
  colors: {
    primaryColor: '#ed646a',
    secondaryColor: '#ed777d',
  },

  /* Custom fonts for website */
  /*
  fonts: {
    myFont: [
      "Times New Roman",
      "Serif"
    ],
    myOtherFont: [
      "-apple-system",
      "system-ui"
    ]
  },
  */

  // This copyright info is used in /core/Footer.js and blog RSS/Atom feeds.
  copyright: `Copyright © ${new Date().getFullYear()} Gojek `,

  highlight: {
    // Highlight.js theme to use for syntax highlighting in code blocks.
    theme: 'default',
  },
  usePrism: ['go', 'bash'],

  // Add custom scripts here that would be placed in <script> tags.
  scripts: ['https://buttons.github.io/buttons.js'],

  // On page navigation for the current documentation page.
  onPageNav: 'separate',
  // No .html extensions for paths.
  cleanUrl: true,

  // For sites with a sizable amount of content, set collapsible to true.
  // Expand/collapse the links and subcategories under categories.
  // docsSideNavCollapsible: true,

  // Show documentation's last contributor's name.
  // enableUpdateBy: true,

  // Show documentation's last update time.
  // enableUpdateTime: true,

  // You may provide arbitrary config keys to be used as needed by your
  // template. For example, if you need your repo's URL...
  //   repoUrl: 'https://github.com/facebook/test-site',
};

module.exports = siteConfig;
