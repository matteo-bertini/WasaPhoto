<script>
export default {
  name: 'CommentItem',
  props: {
    /**
     * The comment object containing ID, author, content, and timestamp.
     */
    comment: { 
      type: Object, 
      required: true 
    },
    /**
     * The username of the post owner to determine moderation privileges.
     */
    postAuthor: { 
      type: String, 
      required: true 
    }
  },
  computed: {
    /**
     * Determines if the current user has permission to remove the comment.
     * Permission is granted to the comment's author or the post's owner.
     */
    canDelete() {
      const loggedUser = localStorage.getItem("Username");
      return loggedUser === this.comment.authorUsername || loggedUser === this.postAuthor;
    },
    /**
     * Formats the ISO timestamp into a human-readable localized string.
     */
    formattedDate() {
      const date = new Date(this.comment.createdAt);
      return date.toLocaleDateString() + " " + date.toLocaleTimeString([], { hour: '2-digit', minute: '2-digit' });
    }
  },
  methods: {
    /**
     * Notifies the parent component to execute the deletion logic 
     * without blocking via confirmation dialogs.
     */
    onDelete() {
      this.$emit('delete-comment', this.comment.commentId);
    }
  }
};
</script>

<template>
  <div class="comment-item">
    <div class="comment-main">
      <div class="comment-header">
        <span class="comment-author">{{ comment.authorUsername }}</span>
        <span class="comment-date">{{ formattedDate }}</span>
      </div>
      
      <div class="comment-body">
        <span class="comment-text">{{ comment.content }}</span>
      </div>
    </div>
    
    <button v-if="canDelete" @click="onDelete" class="delete-comment-btn" title="Delete comment">
      <i class="fa-solid fa-trash-can"></i>
    </button>
  </div>
</template>

<style scoped>
.comment-item {
  display: flex;
  justify-content: space-between;
  align-items: flex-start;
  margin-bottom: 14px;
  gap: 12px;
}

.comment-main {
  display: flex;
  flex-direction: column;
  gap: 4px; 
  flex: 1;
}

.comment-header {
  display: flex;
  align-items: baseline; 
  gap: 8px;
}

.comment-author {
  font-weight: 700;
  font-size: 0.9rem; 
  color: #ffffff;
}

.comment-date {
  font-weight: 400;
  font-size: 0.75rem; 
  color: rgba(255, 255, 255, 0.4);
}

.comment-body {
  line-height: 1.4;
}

.comment-text {
  font-size: 0.85rem; 
  color: #eeeeee;
  word-break: break-word; 
}

.delete-comment-btn {
  background: none;
  border: none;
  color: rgba(255, 255, 255, 0.3);
  cursor: pointer;
  padding: 4px;
  font-size: 0.85rem;
  transition: all 0.2s ease-in-out;
  display: flex;
  align-items: center;
  justify-content: center;
  margin-top: 2px;
}

.delete-comment-btn:hover {
  color: #ff4757;
  transform: scale(1.25);
}
</style>