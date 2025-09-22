import { EventPublisher } from '../../application/ports/EventPublisher.js';

export class LoggingEventPublisher implements EventPublisher {
  constructor(private readonly logger: { info: (message: string, meta?: unknown) => void }) {}

  async publish<T>(topic: string, payload: T): Promise<void> {
    this.logger.info(`Event published to ${topic}`, payload);
  }
}
