<script>
/**
 * @file PostCard.vue
 * @description Component for displaying individual posts, including likes, comments, and media.
 */
import ErrorMsg from './ErrorMsg.vue';
import Comment from './Comment.vue';

export default {
  name: 'PostCard',
  components: { ErrorMsg, Comment },
  props: {
    post: {
      type: Object,
      required: true
    },
  },
  data() {
    return {
      imageUrl: "",
      newCommentText: "",
      likesList: [],
      commentsList: [],
      showLikes: false,
      showComments: false,
      errorMessage: '',
      loading: false,
      // State for the custom deletion modal
      showDeleteConfirm: false 
    };
  },
  computed: {
    /**
     * Identifies if the current user is the author of the post.
     */
    isOwner() {
      return localStorage.getItem("Username") === this.post.username;
    }
  },
  methods: {
    /**
     * Constructs the backend URL to fetch the post image.
     */
    buildImageUrl() {
      const baseUrl = this.$axios.defaults.baseURL;
      this.imageUrl = `${baseUrl}/users/${this.post.username}/posts/${this.post.postId}`;
    },

    /**
     * Handles liking and unliking a post.
     */
    async toggleLike() {
      this.loading = true;
      const token = localStorage.getItem("SessionToken");
      const currentUser = localStorage.getItem("Username");
      try {
        if (this.post.isLikedByMe) {
          await this.$axios.delete(`/users/${this.post.username}/posts/${this.post.postId}/likes`, {
            headers: { Authorization: `Bearer ${token}` }
          });
          this.post.likesNumber--;
          this.likesList = this.likesList.filter(name => name !== currentUser);
        } else {
          await this.$axios.put(`/users/${this.post.username}/posts/${this.post.postId}/likes`, {}, {
            headers: { Authorization: `Bearer ${token}` }
          });
          this.post.likesNumber++;
          if (!this.likesList.includes(currentUser)) this.likesList.push(currentUser);
        }
        this.post.isLikedByMe = !this.post.isLikedByMe;
      } catch (e) {
        this.handleApiError(e);
      } finally {
        this.loading = false;
      }
    },

    /**
     * Retrieves the list of users who liked the post.
     */
    async getLikes() {
      if (this.showLikes) { 
        this.showLikes = false; 
        return; 
      }
      this.loading = true;
      const token = localStorage.getItem("SessionToken");
      try {
        const response = await this.$axios.get(`/users/${this.post.username}/posts/${this.post.postId}/likes`, {
          headers: { Authorization: `Bearer ${token}` }
        });
        this.likesList = response.data;
        this.showLikes = true;
        this.showComments = false;
      } catch (e) {
        this.handleApiError(e);
      } finally {
        this.loading = false;
      }
    },

    /**
     * Retrieves all comments for the current post.
     */
    async getComments() {
      if (this.showComments) { 
        this.showComments = false; 
        return; 
      }
      this.loading = true;
      const token = localStorage.getItem("SessionToken");
      try {
        const response = await this.$axios.get(`/users/${this.post.username}/posts/${this.post.postId}/comments`, {
          headers: { Authorization: `Bearer ${token}` }
        });
        this.commentsList = response.data;
        this.showComments = true;
        this.showLikes = false;
      } catch (e) {
        this.handleApiError(e);
      } finally {
        this.loading = false;
      }
    },

    /**
     * Submits a new comment.
     */
    async addComment() {
      if (!this.newCommentText.trim()) return;
      const token = localStorage.getItem("SessionToken");
      try {
        const response = await this.$axios.post(
          `/users/${this.post.username}/posts/${this.post.postId}/comments`,
          { content: this.newCommentText },
          { headers: { Authorization: `Bearer ${token}` } }
        );
        this.commentsList.push(response.data);
        this.post.commentsNumber++;
        this.newCommentText = "";
      } catch (e) {
        this.handleApiError(e);
      }
    },

    /**
     * Deletes a specific comment from the post.
     */
    async deleteComment(commentId) {
      const token = localStorage.getItem("SessionToken");
      try {
        await this.$axios.delete(
          `/users/${this.post.username}/posts/${this.post.postId}/comments/${commentId}`,
          { headers: { Authorization: `Bearer ${token}` } }
        );
        this.commentsList = this.commentsList.filter(c => c.commentId !== commentId);
        this.post.commentsNumber--;
      } catch (e) {
        this.handleApiError(e);
      }
    },

    /**
     * Permanently removes the post after confirmation.
     */
    async confirmDeletePost() {
      const token = localStorage.getItem("SessionToken");
      try {
        await this.$axios.delete(`/users/${this.post.username}/posts/${this.post.postId}`, {
          headers: { Authorization: `Bearer ${token}` }
        });
        this.$emit('post-deleted', this.post.postId);
        this.showDeleteConfirm = false;
      } catch (e) {
        this.handleApiError(e);
      }
    },

    /**
     * Centralized error handler for the component.
     */
    handleApiError(e) {
      if (e.response) {
        switch (e.response.status) {
          case 401:
            localStorage.clear();
            this.$router.push("/login");
            break;
          case 403:
            this.errorMessage = "Azione non consentita.";
            break;
          case 404:
            this.errorMessage = "Risorsa non trovata.";
            break;
          default:
            this.errorMessage = "Errore del server. Riprova più tardi.";
        }
      } else {
        this.errorMessage = "Connessione al server fallita.";
      }
    }
  },
  mounted() {
    this.buildImageUrl();
  }
};
</script>

<template>
  <div class="post-card glass">
    <ErrorMsg v-if="errorMessage" :message="errorMessage" @close="errorMessage = ''" />

    <div class="post-header">
      <div class="user-info">
        <div class="avatar-mini">{{ post.username.charAt(0).toUpperCase() }}</div>
        <router-link :to="'/users/' + post.username" class="username-link">{{ post.username }}</router-link>
      </div>
      
      <div v-if="post.isSuggested" class="suggested-badge">
        <i class="fa-solid fa-wand-magic-sparkles"></i> Suggerito
      </div>
      
      <button v-if="isOwner" @click="showDeleteConfirm = true" class="btn-icon delete-btn">
        <i class="fa-solid fa-trash-can"></i>
      </button>
    </div>

    <div class="post-media" @dblclick="toggleLike">
      <img v-if="imageUrl" :src="imageUrl" alt="Post photo" />
      <div v-else class="image-loader"><i class="fa-solid fa-circle-notch fa-spin"></i></div>
    </div>

    <div class="post-footer">
      <div class="actions-row">
        <div class="action-group">
          <button @click="toggleLike" class="btn-icon heart-icon" :class="{ 'liked': post.isLikedByMe }">
            <i :class="post.isLikedByMe ? 'fa-solid fa-heart' : 'fa-regular fa-heart'"></i>
          </button>
          <span class="stat-text" @click="getLikes">{{ post.likesNumber }}</span>
        </div>

        <div class="action-group">
          <button @click="getComments" class="btn-icon comment-icon"><i class="fa-regular fa-comment"></i></button>
          <span class="stat-text" @click="getComments">{{ post.commentsNumber }}</span>
        </div>
      </div>

      <div v-if="post.caption" class="post-caption">
        <span class="caption-username">{{ post.username }}</span> {{ post.caption }}
      </div>

      <!-- Likes List Expansion -->
      <div v-if="showLikes" class="expanded-list">
        <p class="list-title">Piace a:</p>
        <div class="list-content">
          {{ likesList.length > 0 ? likesList.join(', ') : 'Ancora nessun mi piace.' }}
        </div>
      </div>

      <!-- Comments List Expansion -->
      <div v-if="showComments" class="expanded-list">
        <p class="list-title">Commenti</p>
        <div class="list-content scrollable">
          <Comment 
            v-for="c in commentsList" 
            :key="c.commentId" 
            :comment="c" 
            :postAuthor="post.username"
            @delete-comment="deleteComment"
          />
          <p v-if="commentsList.length === 0" class="empty-msg">Ancora nessun commento.</p>
        </div>

        <div class="add-comment-area">
          <input 
            type="text" v-model="newCommentText" 
            placeholder="Aggiungi un commento..." 
            @keyup.enter="addComment" class="comment-input"
          />
          <button class="btn-post-comment" :disabled="!newCommentText.trim()" @click="addComment">Invia</button>
        </div>
      </div>
    </div>

    <!-- Custom Delete Confirmation Modal -->
    <div v-if="showDeleteConfirm" class="modal-overlay" @click.self="showDeleteConfirm = false">
      <div class="glass-card confirm-modal">
        <h3>Elimina Post</h3>
        <p>Sei sicuro di voler eliminare definitivamente questo post? L'azione non è reversibile.</p>
        <div class="modal-actions">
          <button class="btn-cancel" @click="showDeleteConfirm = false">Annulla</button>
          <button class="gradient-button btn-delete-final" @click="confirmDeletePost">Elimina</button>
        </div>
      </div>
    </div>
  </div>
</template>

<style scoped>
.post-card {
  max-width: 500px; background: rgba(255, 255, 255, 0.05);
  backdrop-filter: blur(10px); border: 1px solid rgba(255, 255, 255, 0.1);
  border-radius: 12px; margin: 0 auto 30px; color: white;
  box-shadow: 0 10px 30px rgba(0,0,0,0.5); overflow: hidden;
}

.post-header { padding: 12px 15px; display: flex; justify-content: space-between; align-items: center; }
.user-info { display: flex; align-items: center; gap: 10px; }
.avatar-mini {
  width: 32px; height: 32px; border-radius: 50%;
  background: linear-gradient(135deg,#003366 100%); color: white;
  display: flex; align-items: center; justify-content: center; font-weight: bold; font-size: 0.8rem;
}
.username-link { color: white; text-decoration: none; font-weight: 600; font-size: 0.9rem; }

.suggested-badge {
  font-size: 0.65rem; color: #003366; background: #fff;
  padding: 4px 10px; border-radius: 20px; font-weight: bold; text-transform: uppercase;
}

.post-media { width: 100%; aspect-ratio: 1/1; background: #000; display: flex; align-items: center; justify-content: center; }
.post-media img { width: 100%; height: 100%; object-fit: cover; }

.post-footer { padding: 12px 15px; }
.actions-row { display: flex; gap: 15px; margin-bottom: 10px; }
.action-group { display: flex; align-items: center; gap: 6px; }
.btn-icon { background: none; border: none; color: white; cursor: pointer; font-size: 1.3rem; transition: 0.2s; }
.btn-icon:hover { transform: scale(1.1); }
.liked i { color: #66001a; }
.stat-text { font-size: 0.85rem; font-weight: 600; cursor: pointer; }
.delete-btn { color: rgba(255,255,255,0.4); font-size: 1rem; }

.post-caption { font-size: 0.9rem; margin-top: 5px; }
.caption-username { font-weight: 700; }

.expanded-list { margin-top: 15px; padding-top: 10px; border-top: 1px solid rgba(255,255,255,0.1); }
.list-title { font-weight: 700; color: #aaa; margin-bottom: 8px; font-size: 0.8rem; }
.scrollable { max-height: 180px; overflow-y: auto; }

.add-comment-area {
  display: flex; align-items: center; background: rgba(255,255,255,0.05);
  border-radius: 8px; padding: 6px 12px; border: 1px solid rgba(255,255,255,0.1); margin-top: 10px;
}
.comment-input { flex: 1; background: transparent; border: none; color: white; font-size: 0.85rem; outline: none; }
.btn-post-comment { background: none; border: none; color: #0095f6; font-weight: 700; cursor: pointer; }

/* Custom Modal Styles */
.modal-overlay {
  position: fixed; top: 0; left: 0; width: 100%; height: 100%;
  background: rgba(0, 0, 0, 0.8); backdrop-filter: blur(5px);
  display: flex; align-items: center; justify-content: center; z-index: 5000;
}
.glass-card {
  background: rgba(20, 20, 20, 0.95); border: 1px solid rgba(255, 255, 255, 0.1);
  border-radius: 20px; padding: 25px; text-align: center; max-width: 400px;
}
.confirm-modal h3 { color: #fff; margin-bottom: 15px; }
.confirm-modal p { color: #ccc; font-size: 0.9rem; margin-bottom: 25px; }
.modal-actions { display: flex; justify-content: space-between; gap: 20px; }
.btn-cancel {
  background: rgba(255, 255, 255, 0.1); color: white; border: none;
  padding: 10px 20px; border-radius: 10px; cursor: pointer; flex: 1;
}
.btn-delete-final { background: #66001a !important; flex: 1; border: none; border-radius: 10px; color: white; cursor: pointer; }
.gradient-button { font-weight: 600; padding: 10px 20px; }
</style>