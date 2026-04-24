<script>
import ErrorMsg from './ErrorMsg.vue';

export default {
  name: 'PostCard',
  components: { ErrorMsg },
  props: {
    post: {
      type: Object,
      required: true
    }
  },
  data() {
    return {
      likesList: [],
      commentsList: [],
      showLikes: false,
      showComments: false,
      errorMsg: '',
      loading: false
    };
  },
  computed: {
    // Check if the current user is the owner of the post
    isOwner() {
      return localStorage.getItem("Username") === this.post.username;
    }
  },
  methods: {
    setPostImage() {
      // We construct the URL directly pointing to the backend binary endpoint
      const baseUrl = this.$axios.defaults.baseURL;
      this.imageUrl = `${baseUrl}/users/${this.post.username}/posts/${this.post.postId}`;
    },

    async toggleLike() {
      const token = localStorage.getItem("SessionToken");
      try {
        if (this.post.isLikedByMe) {
          await this.$axios.delete(`/users/${this.post.username}/posts/${this.post.postId}/likes`, {
            headers: { Authorization: `Bearer ${token}` }
          });
          this.post.likesNumber--;
        } else {
          await this.$axios.put(`/users/${this.post.username}/posts/${this.post.postId}/likes`, {}, {
            headers: { Authorization: `Bearer ${token}` }
          });
          this.post.likesNumber++;
        }
        this.post.isLikedByMe = !this.post.isLikedByMe;
      } catch (e) {
        this.handleError(e, "Could not update like");
      }
    },

    async fetchLikes() {
      if (this.showLikes) { this.showLikes = false; return; }
      const token = localStorage.getItem("SessionToken");
      try {
        const response = await this.$axios.get(`/users/${this.post.username}/posts/${this.post.postId}/likes`, {
          headers: { Authorization: `Bearer ${token}` }
        });
        this.likesList = response.data;
        this.showLikes = true;
        this.showComments = false;
      } catch (e) {
        this.handleError(e, "Error fetching likes");
      }
    },

    async fetchComments() {
      if (this.showComments) { this.showComments = false; return; }
      const token = localStorage.getItem("SessionToken");
      try {
        const response = await this.$axios.get(`/users/${this.post.username}/posts/${this.post.postId}/comments`, {
          headers: { Authorization: `Bearer ${token}` }
        });
        this.commentsList = response.data;
        this.showComments = true;
        this.showLikes = false;
      } catch (e) {
        this.handleError(e, "Error fetching comments");
      }
    },
    async deletePost() {
      if (!confirm("Are you sure you want to delete this post?")) return;
      const token = localStorage.getItem("SessionToken");
      try {
        await this.$axios.delete(`/users/${this.post.username}/posts/${this.post.postId}`, {
          headers: { Authorization: `Bearer ${token}` }
        });
        // Emit to parent (ProfileView) to update the list locally
        this.$emit('post-deleted', this.post.postId);
      } catch (e) {
        this.handleError(e, "Failed to delete post");
      }
    },

    handleError(e, defaultMsg) {
      console.error(`[PostCard Error] ${defaultMsg}:`, e);
      if (e.response && e.response.status === 403) {
        this.errorMsg = "Forbidden: Access blocked by a ban.";
      } else {
        this.errorMsg = defaultMsg;
      }
    }
  },
  mounted() {
    this.setPostImage();
  }
};
</script>
<template>
  <div class="post-card glass">
    <ErrorMsg v-if="errorMsg" :message="errorMsg" @close="errorMsg = ''" />

    <div class="post-header">
      <div class="user-info">
        <div class="avatar-mini">{{ post.username.charAt(0).toUpperCase() }}</div>
        <router-link :to="'/profiles/' + post.username" class="username-link">
          {{ post.username }}
        </router-link>
      </div>
      
      <button v-if="isOwner" @click="deletePost" class="btn-icon delete-btn" title="Delete Post">
        <i class="fa-solid fa-trash-can"></i>
      </button>
    </div>

    <div class="post-media" @dblclick="toggleLike">
      <img v-if="imageUrl" :src="imageUrl" alt="Post content" />
      <div v-else class="image-loader">
        <i class="fa-solid fa-circle-notch fa-spin"></i>
      </div>
    </div>

    <div class="post-footer">
      <div class="actions-row">
        <div class="action-group">
          <button @click="toggleLike" class="btn-icon heart-icon" :class="{ 'liked': post.isLikedByMe }">
            <i :class="post.isLikedByMe ? 'fa-solid fa-heart' : 'fa-regular fa-heart'"></i>
          </button>
          <span class="stat-text" @click="fetchLikes">{{ post.likesCount }} likes</span>
        </div>

        <div class="action-group">
          <button @click="fetchComments" class="btn-icon comment-icon">
            <i class="fa-regular fa-comment"></i>
          </button>
          <span class="stat-text" @click="fetchComments">{{ post.commentsCount }} comments</span>
        </div>
      </div>

      <div class="post-caption">
        <span class="caption-username">{{ post.username }}</span> {{ post.caption }}
      </div>

      <div v-if="showLikes" class="expanded-list">
        <p class="list-title">Liked by:</p>
        <div class="list-content">
          {{ likesList.length > 0 ? likesList.join(', ') : 'No likes yet' }}
        </div>
      </div>

      <div v-if="showComments" class="expanded-list">
        <p class="list-title">Comments:</p>
        <div class="list-content">
          <div v-for="c in commentsList" :key="c.commentId" class="comment-item">
            <strong>{{ c.username }}</strong> {{ c.content }}
          </div>
          <p v-if="commentsList.length === 0">No comments yet.</p>
        </div>
      </div>
    </div>
  </div>
</template>



<style scoped>
.post-card {
  max-width: 500px;
  background: rgba(255, 255, 255, 0.05);
  backdrop-filter: blur(10px);
  border: 1px solid rgba(255, 255, 255, 0.1);
  border-radius: 12px;
  margin: 0 auto 30px;
  overflow: hidden;
  color: white;
  box-shadow: 0 10px 30px rgba(0,0,0,0.5);
}

.post-header {
  padding: 12px 15px;
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.user-info {
  display: flex;
  align-items: center;
  gap: 10px;
}

.avatar-mini {
  width: 32px; height: 32px; border-radius: 50%;
  background: linear-gradient(45deg, #f09433, #e6683c, #dc2743, #cc2366, #bc1888);
  display: flex; align-items: center; justify-content: center;
  font-weight: bold; font-size: 0.8rem;
}

.username-link {
  color: white; text-decoration: none; font-weight: 600; font-size: 0.9rem;
}

.post-media {
  width: 100%;
  aspect-ratio: 1 / 1;
  background: #000;
  display: flex; align-items: center; justify-content: center;
  overflow: hidden;
}

.post-media img {
  width: 100%; height: 100%; object-fit: cover;
}

.post-footer {
  padding: 12px 15px;
}

.actions-row {
  display: flex; gap: 15px; margin-bottom: 10px;
}

.action-group {
  display: flex; align-items: center; gap: 6px;
}

.btn-icon {
  background: none; border: none; color: white; cursor: pointer;
  font-size: 1.4rem; padding: 0; transition: transform 0.1s ease;
}

.btn-icon:hover { transform: scale(1.1); }
.liked i { color: #ff4757; }

.stat-text {
  font-size: 0.85rem; font-weight: 600; cursor: pointer;
}

.delete-btn {
  font-size: 1.1rem; color: rgba(255, 255, 255, 0.4);
}
.delete-btn:hover { color: #ff4757; }

.post-caption {
  font-size: 0.9rem; line-height: 1.4; margin-top: 5px;
}

.caption-username { font-weight: 700; margin-right: 5px; }

.expanded-list {
  margin-top: 15px; padding-top: 10px;
  border-top: 1px solid rgba(255, 255, 255, 0.1);
  font-size: 0.85rem;
}

.list-title { font-weight: 700; color: #aaa; margin-bottom: 5px; }

.comment-item { margin-bottom: 4px; }

.image-loader { font-size: 2rem; color: rgba(255,255,255,0.2); }
</style>