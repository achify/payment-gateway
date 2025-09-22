import { createServer } from 'http';
import { createApp } from './config/app.js';

const PORT = process.env.PORT ?? 3000;

const app = createApp();

createServer(app).listen(PORT, () => {
  console.log(`Payment gateway API listening on port ${PORT}`);
});
